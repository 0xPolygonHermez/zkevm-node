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
	DataCommitteeContract *cdkdatacommittee.Cdkdatacommittee
	isTrustedSequencer    bool
	l2Coinbase            common.Address

	privKey                    *ecdsa.PrivateKey
	state                      stateInterface
	zkEVMClient                syncinterfaces.ZKEVMClientInterface
	dataCommitteeClientFactory client.ClientFactoryInterface

	ctx                     context.Context
	committeeMembers        []DataCommitteeMember
	selectedCommitteeMember int
}

type Config struct {
	L1RPCURL           string
	DataCommitteeAddr  common.Address
	IsTrustedSequencer bool
	L2Coinbase         common.Address
}

func NewDataCommitteeMan(
	c Config,
	privKey *ecdsa.PrivateKey,
	state stateInterface,
	zkEVMClient syncinterfaces.ZKEVMClientInterface,
	dataCommitteeClientFactory client.ClientFactoryInterface,
) (*DataCommitteeMan, error) {
	ethClient, err := ethclient.Dial(c.L1RPCURL)
	if err != nil {
		log.Errorf("error connecting to %s: %+v", c.L1RPCURL, err)
		return nil, err
	}
	dataCommittee, err := cdkdatacommittee.NewCdkdatacommittee(c.DataCommitteeAddr, ethClient)
	if err != nil {
		return nil, err
	}
	dacman := &DataCommitteeMan{
		DataCommitteeContract:      dataCommittee,
		isTrustedSequencer:         c.IsTrustedSequencer,
		l2Coinbase:                 c.L2Coinbase,
		privKey:                    privKey,
		state:                      state,
		zkEVMClient:                zkEVMClient,
		dataCommitteeClientFactory: dataCommitteeClientFactory,
		ctx:                        context.Background(),
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
