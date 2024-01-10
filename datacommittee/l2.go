package datacommittee

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"strings"

	"github.com/0xPolygon/cdk-data-availability/client"
	jTypes "github.com/0xPolygon/cdk-data-availability/rpc"
	daTypes "github.com/0xPolygon/cdk-data-availability/types"
	"github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const unexpectedHashTemplate = "missmatch on transaction data for batch num %d. Expected hash %s, actual hash: %s"

func (d *DataCommitteeMan) GetBatchL2Data(batchNum uint64, expectedTransactionsHash common.Hash) ([]byte, error) {
	found := true
	transactionsData, err := d.state.GetBatchL2DataByNumber(d.ctx, batchNum, nil)
	if err != nil {
		if err == state.ErrNotFound {
			found = false
		} else {
			return nil, fmt.Errorf("failed to get batch data from state for batch num %d: %w", batchNum, err)
		}
	}
	actualTransactionsHash := crypto.Keccak256Hash(transactionsData)
	if !found || expectedTransactionsHash != actualTransactionsHash {
		if found {
			log.Warnf(unexpectedHashTemplate, batchNum, expectedTransactionsHash, actualTransactionsHash)
		}

		if !d.isTrustedSequencer {
			log.Info("trying to get data from trusted sequencer")
			data, err := d.getDataFromTrustedSequencer(batchNum, expectedTransactionsHash)
			if err != nil {
				log.Error(err)
			} else {
				return data, nil
			}
		}

		log.Info("trying to get data from data committee node")
		data, err := d.getDataFromCommittee(batchNum, expectedTransactionsHash)
		if err != nil {
			log.Error(err)
			if d.isTrustedSequencer {
				return nil, fmt.Errorf("data not found on the local DB nor on any data committee member")
			} else {
				return nil, fmt.Errorf("data not found on the local DB, nor from the trusted sequencer nor on any data committee member")
			}
		}
		return data, nil
	}
	return transactionsData, nil
}

func (d *DataCommitteeMan) getDataFromCommittee(batchNum uint64, expectedTransactionsHash common.Hash) ([]byte, error) {
	intialMember := d.selectedCommitteeMember
	found := false
	for !found && intialMember != -1 {
		member := d.committeeMembers[d.selectedCommitteeMember]
		log.Infof("trying to get data from %s at %s", member.Addr.Hex(), member.URL)
		c := d.dataCommitteeClientFactory.New(member.URL)
		data, err := c.GetOffChainData(d.ctx, expectedTransactionsHash)
		if err != nil {
			log.Warnf(
				"error getting data from DAC node %s at %s: %s",
				member.Addr.Hex(), member.URL, err,
			)
			d.selectedCommitteeMember = (d.selectedCommitteeMember + 1) % len(d.committeeMembers)
			if d.selectedCommitteeMember == intialMember {
				break
			}
			continue
		}
		actualTransactionsHash := crypto.Keccak256Hash(data)
		if actualTransactionsHash != expectedTransactionsHash {
			unexpectedHash := fmt.Errorf(
				unexpectedHashTemplate, batchNum, expectedTransactionsHash, actualTransactionsHash,
			)
			log.Warnf(
				"error getting data from DAC node %s at %s: %s",
				member.Addr.Hex(), member.URL, unexpectedHash,
			)
			d.selectedCommitteeMember = (d.selectedCommitteeMember + 1) % len(d.committeeMembers)
			if d.selectedCommitteeMember == intialMember {
				break
			}
			continue
		}
		return data, nil
	}
	if err := d.loadCommittee(); err != nil {
		return nil, fmt.Errorf("error loading data committee: %s", err)
	}
	return nil, fmt.Errorf("couldn't get the data from any committee member")
}

func (d *DataCommitteeMan) getDataFromTrustedSequencer(batchNum uint64, expectedTransactionsHash common.Hash) ([]byte, error) {
	b, err := d.zkEVMClient.BatchByNumber(d.ctx, big.NewInt(int64(batchNum)))
	if err != nil {
		return nil, fmt.Errorf("failed to get batch num %d from trusted sequencer: %w", batchNum, err)
	}
	actualTransactionsHash := crypto.Keccak256Hash(b.BatchL2Data)
	if expectedTransactionsHash != actualTransactionsHash {
		return nil, fmt.Errorf(
			unexpectedHashTemplate, batchNum, expectedTransactionsHash, actualTransactionsHash,
		)
	}
	return b.BatchL2Data, nil
}

type signatureMsg struct {
	addr      common.Address
	signature []byte
	err       error
}

func (s *DataCommitteeMan) GetSignaturesAndAddrsFromDataCommittee(ctx context.Context, sequences []types.Sequence) ([]byte, error) {
	// Get current committee
	committee, err := s.GetCurrentDataCommittee()
	if err != nil {
		return nil, err
	}

	// Get last accInputHash
	var accInputHash common.Hash
	if sequences[0].BatchNumber != 0 {
		prevBatch, err := s.state.GetBatchByNumber(ctx, sequences[0].BatchNumber-1, nil)
		if err != nil {
			return nil, err
		}
		accInputHash = prevBatch.AccInputHash
	}

	// Authenticate as trusted sequencer by signing the sequences
	sequence := daTypes.Sequence{
		Batches:         []daTypes.Batch{},
		OldAccInputHash: accInputHash,
	}
	for _, seq := range sequences {
		sequence.Batches = append(sequence.Batches, daTypes.Batch{
			Number:         jTypes.ArgUint64(seq.BatchNumber),
			GlobalExitRoot: seq.GlobalExitRoot,
			Timestamp:      jTypes.ArgUint64(seq.Timestamp),
			Coinbase:       s.l2Coinbase,
			L2Data:         seq.BatchL2Data,
		})
	}
	signedSequence, err := sequence.Sign(s.privKey)
	if err != nil {
		return nil, err
	}

	// Request signatures to all members in parallel
	ch := make(chan signatureMsg, len(committee.Members))
	signatureCtx, cancelSignatureCollection := context.WithCancel(ctx)
	for _, member := range committee.Members {
		go requestSignatureFromMember(signatureCtx, *signedSequence, member, ch)
	}

	// Collect signatures
	msgs := []signatureMsg{}
	var collectedSignatures uint64
	var failedToCollect uint64
	for collectedSignatures < committee.RequiredSignatures {
		msg := <-ch
		if msg.err != nil {
			log.Errorf("error when trying to get signature from %s: %s", msg.addr, msg.err)
			failedToCollect++
			if len(committee.Members)-int(failedToCollect) < int(committee.RequiredSignatures) {
				cancelSignatureCollection()
				return nil, errors.New("too many members failed to send their signature")
			}
		} else {
			log.Infof("received signature from %s", msg.addr)
			collectedSignatures++
		}
		msgs = append(msgs, msg)
	}

	// Stop requesting as soon as we have N valid signatures
	cancelSignatureCollection()

	return buildSignaturesAndAddrs(signatureMsgs(msgs), committee.Members), nil
}

func requestSignatureFromMember(ctx context.Context, signedSequence daTypes.SignedSequence, member DataCommitteeMember, ch chan signatureMsg) {
	// request
	c := client.New(member.URL)
	log.Infof("sending request to sign the sequence to %s at %s", member.Addr.Hex(), member.URL)
	signature, err := c.SignSequence(signedSequence)
	if err != nil {
		ch <- signatureMsg{
			addr: member.Addr,
			err:  err,
		}
		return
	}
	// verify returned signature
	signedSequence.Signature = signature
	signer, err := signedSequence.Signer()
	if err != nil {
		ch <- signatureMsg{
			addr: member.Addr,
			err:  err,
		}
		return
	}
	if signer != member.Addr {
		ch <- signatureMsg{
			addr: member.Addr,
			err:  fmt.Errorf("invalid signer. Expected %s, actual %s", member.Addr.Hex(), signer.Hex()),
		}
		return
	}
	ch <- signatureMsg{
		addr:      member.Addr,
		signature: signature,
	}
}

func buildSignaturesAndAddrs(msgs signatureMsgs, members []DataCommitteeMember) []byte {
	res := []byte{}
	sort.Sort(msgs)
	for _, msg := range msgs {
		log.Debugf("adding signature %s from %s", common.Bytes2Hex(msg.signature), msg.addr.Hex())
		res = append(res, msg.signature...)
	}
	for _, member := range members {
		log.Debugf("adding addr %s", common.Bytes2Hex(member.Addr.Bytes()))
		res = append(res, member.Addr.Bytes()...)
	}
	log.Debugf("full res %s", common.Bytes2Hex(res))
	return res
}

type signatureMsgs []signatureMsg

func (s signatureMsgs) Len() int { return len(s) }
func (s signatureMsgs) Less(i, j int) bool {
	return strings.ToUpper(s[i].addr.Hex()) < strings.ToUpper(s[j].addr.Hex())
}
func (s signatureMsgs) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
