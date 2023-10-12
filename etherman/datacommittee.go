package etherman

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
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

// GetCurrentDataCommittee return the currently registered data committee
func (etherMan *Client) GetCurrentDataCommittee() (*DataCommittee, error) {
	addrsHash, err := etherMan.DataCommittee.CommitteeHash(&bind.CallOpts{Pending: false})
	if err != nil {
		return nil, fmt.Errorf("error getting CommitteeHash from L1 SC: %w", err)
	}
	reqSign, err := etherMan.DataCommittee.RequiredAmountOfSignatures(&bind.CallOpts{Pending: false})
	if err != nil {
		return nil, fmt.Errorf("error getting RequiredAmountOfSignatures from L1 SC: %w", err)
	}
	members, err := etherMan.GetCurrentDataCommitteeMembers()
	if err != nil {
		return nil, err
	}

	return &DataCommittee{
		AddressesHash:      common.Hash(addrsHash),
		RequiredSignatures: reqSign.Uint64(),
		Members:            members,
	}, nil
}

// GetCurrentDataCommitteeMembers return the currently registered data committee members
func (etherMan *Client) GetCurrentDataCommitteeMembers() ([]DataCommitteeMember, error) {
	members := []DataCommitteeMember{}
	nMembers, err := etherMan.DataCommittee.GetAmountOfMembers(&bind.CallOpts{Pending: false})
	if err != nil {
		return nil, fmt.Errorf("error getting GetAmountOfMembers from L1 SC: %w", err)
	}
	for i := int64(0); i < nMembers.Int64(); i++ {
		member, err := etherMan.DataCommittee.Members(&bind.CallOpts{Pending: false}, big.NewInt(i))
		if err != nil {
			return nil, fmt.Errorf("error getting Members %d from L1 SC: %w", i, err)
		}
		members = append(members, DataCommitteeMember{
			Addr: member.Addr,
			URL:  member.Url,
		})
	}
	return members, nil
}
