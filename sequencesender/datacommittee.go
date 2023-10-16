package sequencesender

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/0xPolygon/cdk-data-availability/batch"
	"github.com/0xPolygon/cdk-data-availability/client"
	"github.com/0xPolygon/cdk-data-availability/rpc"
	"github.com/0xPolygon/cdk-data-availability/sequence"
	ethman "github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/common"
)

type signatureMsg struct {
	addr      common.Address
	signature []byte
	err       error
}

func (s *SequenceSender) getSignaturesAndAddrsFromDataCommittee(ctx context.Context, sequences []types.Sequence) ([]byte, error) {
	// Get current committee
	committee, err := s.etherman.GetCurrentDataCommittee()
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
	sequence := sequence.Sequence{
		Batches:         []batch.Batch{},
		OldAccInputHash: accInputHash,
	}
	for _, seq := range sequences {
		sequence.Batches = append(sequence.Batches, batch.Batch{
			Number:         rpc.ArgUint64(seq.BatchNumber),
			GlobalExitRoot: seq.GlobalExitRoot,
			Timestamp:      rpc.ArgUint64(seq.Timestamp),
			Coinbase:       s.cfg.L2Coinbase,
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

func requestSignatureFromMember(ctx context.Context, signedSequence sequence.SignedSequence, member ethman.DataCommitteeMember, ch chan signatureMsg) {
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

func buildSignaturesAndAddrs(msgs signatureMsgs, members []ethman.DataCommitteeMember) []byte {
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
