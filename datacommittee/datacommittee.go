package datacommittee

import (
	"context"
	"crypto/ecdsa"
	"math/rand"

	"github.com/0xPolygon/cdk-data-availability/client"
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/cdkdatacommittee"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/common/syncinterfaces"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// DataCommitteeMember represents a member of the Data Committee
type DataCommitteeMember struct {
	Addr common.Address
	URL  string
}

// DataCommittee represents a specific committee
type DataCommittee struct {
	AddressesHash      common.Hash
	Members            []DataCommitteeMember
	RequiredSignatures uint64
}

type DataCommitteeMan struct {
	DataCommitteeContract      *cdkdatacommittee.Cdkdatacommittee
	isTrustedSequencer         bool
	zkEVMClient                syncinterfaces.ZKEVMClientInterface
	committeeMembers           []DataCommitteeMember
	state                      stateInterface
	selectedCommitteeMember    int
	dataCommitteeClientFactory client.ClientFactoryInterface
	ctx                        context.Context
	l2Coinbase                 common.Address
	privKey                    *ecdsa.PrivateKey
}

func NewDataCommitteeMan(dataCommitteeAddr common.Address, l1RPCURL string) (*DataCommitteeMan, error) {
	ethClient, err := ethclient.Dial(l1RPCURL)
	if err != nil {
		log.Errorf("error connecting to %s: %+v", l1RPCURL, err)
		return nil, err
	}
	dataCommittee, err := cdkdatacommittee.NewCdkdatacommittee(dataCommitteeAddr, ethClient)
	if err != nil {
		return nil, err
	}
	dacman := &DataCommitteeMan{
		DataCommitteeContract: dataCommittee,
	}
	err = dacman.loadCommittee()
	return dacman, err
}

func (d *DataCommitteeMan) loadCommittee() error {
	committee, err := d.GetCurrentDataCommittee()
	if err != nil {
		return err
	}
	selectedCommitteeMember := -1
	if committee != nil {
		d.committeeMembers = committee.Members
		if len(committee.Members) > 0 {
			selectedCommitteeMember = rand.Intn(len(committee.Members)) //nolint:gosec
		}
	}
	d.selectedCommitteeMember = selectedCommitteeMember
	return nil
}
