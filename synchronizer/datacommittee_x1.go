package synchronizer

import (
	"fmt"
	"math/big"
	"math/rand"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const unexpectedHashTemplate = "missmatch on transaction data for batch num %d. Expected hash %s, actual hash: %s"

func (s *ClientSynchronizer) loadCommittee() error {
	committee, err := s.etherMan.GetCurrentDataCommittee()
	if err != nil {
		return err
	}
	selectedCommitteeMember := -1
	if committee != nil {
		s.committeeMembers = committee.Members
		if len(committee.Members) > 0 {
			selectedCommitteeMember = rand.Intn(len(committee.Members)) //nolint:gosec
		}
	}
	s.selectedCommitteeMember = selectedCommitteeMember
	return nil
}

func (s *ClientSynchronizer) getBatchL2Data(batchNum uint64, expectedTransactionsHash common.Hash) ([]byte, error) {
	found := true
	transactionsData, err := s.state.GetBatchL2DataByNumber(s.ctx, batchNum, nil)
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

		if !s.isTrustedSequencer {
			log.Info("trying to get data from trusted sequencer")
			data, err := s.getDataFromTrustedSequencer(batchNum, expectedTransactionsHash)
			if err != nil {
				log.Error(err)
			} else {
				return data, nil
			}
		}

		log.Info("trying to get data from data committee node")
		data, err := s.getDataFromCommittee(batchNum, expectedTransactionsHash)
		if err != nil {
			log.Error(err)
			if s.isTrustedSequencer {
				return nil, fmt.Errorf("data not found on the local DB nor on any data committee member")
			} else {
				return nil, fmt.Errorf("data not found on the local DB, nor from the trusted sequencer nor on any data committee member")
			}
		}
		return data, nil
	}
	return transactionsData, nil
}

func (s *ClientSynchronizer) getDataFromCommittee(batchNum uint64, expectedTransactionsHash common.Hash) ([]byte, error) {
	intialMember := s.selectedCommitteeMember
	found := false
	for !found && intialMember != -1 {
		member := s.committeeMembers[s.selectedCommitteeMember]
		log.Infof("trying to get data from %s at %s", member.Addr.Hex(), member.URL)
		c := s.dataCommitteeClientFactory.New(member.URL)
		data, err := c.GetOffChainData(s.ctx, expectedTransactionsHash)
		if err != nil {
			log.Warnf(
				"error getting data from DAC node %s at %s: %s",
				member.Addr.Hex(), member.URL, err,
			)
			s.selectedCommitteeMember = (s.selectedCommitteeMember + 1) % len(s.committeeMembers)
			if s.selectedCommitteeMember == intialMember {
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
			s.selectedCommitteeMember = (s.selectedCommitteeMember + 1) % len(s.committeeMembers)
			if s.selectedCommitteeMember == intialMember {
				break
			}
			continue
		}
		return data, nil
	}
	if err := s.loadCommittee(); err != nil {
		return nil, fmt.Errorf("error loading data committee: %s", err)
	}
	return nil, fmt.Errorf("couldn't get the data from any committee member")
}

func (s *ClientSynchronizer) getDataFromTrustedSequencer(batchNum uint64, expectedTransactionsHash common.Hash) ([]byte, error) {
	b, err := s.zkEVMClient.BatchByNumber(s.ctx, big.NewInt(int64(batchNum)))
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
