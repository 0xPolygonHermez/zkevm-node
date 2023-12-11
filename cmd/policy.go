package main

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/0xPolygonHermez/zkevm-node/config"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/pool/pgpoolstorage"
	"github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli/v2"
)

var (
	policyFlag = cli.StringFlag{
		Name:     "policy",
		Aliases:  []string{"p"},
		Usage:    "Name of policy to operate on",
		Required: false,
	}
	csvFlag = cli.StringFlag{
		Name:     "csv",
		Usage:    "CSV file with addresses",
		Required: false,
	}
	allowFlag = cli.BoolFlag{
		Name:     "allow",
		Usage:    "Update policy to 'allow' addresses on list",
		Required: false,
	}
	denyFlag = cli.BoolFlag{
		Name:     "deny",
		Usage:    "Update policy to 'deny' addresses on list",
		Required: false,
	}
	noHeaderFlag = cli.BoolFlag{
		Name:     "no-header",
		Value:    false,
		Required: false,
	}

	policyActionFlags = []cli.Flag{&policyFlag}
)

var policyCommands = cli.Command{
	Name:   "policy",
	Usage:  "View, update, and apply policies",
	Action: describe,
	Flags:  []cli.Flag{&configFileFlag},
	Subcommands: []*cli.Command{
		{
			Name:   "add",
			Usage:  "Add address(es) to a policy exclusion list",
			Action: addAcl,
			Flags:  append(policyActionFlags, &csvFlag),
		}, {
			Name:   "clear",
			Usage:  "Clear the addresses listed as exceptions to a policy",
			Action: clearAcl,
			Flags:  policyActionFlags,
		}, {
			Name:   "describe",
			Usage:  "Describe the default actions for the policies",
			Action: describe,
			Flags:  append(policyActionFlags, &noHeaderFlag),
		}, {
			Name:   "remove",
			Usage:  "Remove address(es) from a policy exclusion list",
			Action: removeAcl,
			Flags:  append(policyActionFlags, &csvFlag),
		}, {
			Name:   "update",
			Usage:  "Update the default action for a policy",
			Action: updatePolicy,
			Flags:  append(policyActionFlags, &allowFlag, &denyFlag),
		},
	},
}

func updatePolicy(cli *cli.Context) error {
	_, db, err := configAndStorage(cli)
	if err != nil {
		return err
	}
	policy, err := resolvePolicy(cli)
	if err != nil {
		return err
	}

	allow := cli.Bool(allowFlag.Name)
	deny := cli.Bool(denyFlag.Name)

	// exactly one must be set
	if (allow && deny) || (!allow && !deny) {
		return errors.New("supply one policy action [--allow or --deny]")
	}

	var setting bool
	if allow {
		setting = true
	} else if deny {
		setting = false
	}

	err = db.UpdatePolicy(context.Background(), policy, setting)
	if err != nil {
		return err
	}
	return nil
}

func addAcl(cli *cli.Context) error {
	_, db, err := configAndStorage(cli)
	if err != nil {
		return err
	}
	policy, addresses, err := requirePolicyAndAddresses(cli)
	if err != nil {
		return err
	}
	err = db.AddAddressesToPolicy(context.Background(), policy, addresses)
	if err != nil {
		return err
	}
	return nil
}

func removeAcl(cli *cli.Context) error {
	_, db, err := configAndStorage(cli)
	if err != nil {
		return err
	}
	policy, addresses, err := requirePolicyAndAddresses(cli)
	if err != nil {
		return err
	}
	err = db.RemoveAddressesFromPolicy(context.Background(), policy, addresses)
	if err != nil {
		return err
	}
	return nil
}

func clearAcl(cli *cli.Context) error {
	_, db, err := configAndStorage(cli)
	if err != nil {
		return err
	}
	policy, err := resolvePolicy(cli)
	if err != nil {
		return err
	}
	err = db.ClearPolicy(context.Background(), policy)
	if err != nil {
		return err
	}
	return nil
}

func describe(cli *cli.Context) error {
	showHeader := !cli.Bool(noHeaderFlag.Name)
	if cli.IsSet(policyFlag.Name) {
		return describePolicy(cli, showHeader)
	}
	return describePolicies(cli, showHeader)
}

func describePolicy(cli *cli.Context, showHeader bool) error {
	_, db, err := configAndStorage(cli)
	if err != nil {
		return err
	}

	policyName, err := resolvePolicy(cli)
	if err != nil {
		return err
	}

	if showHeader {
		policy, err := db.DescribePolicy(context.Background(), policyName)
		if err != nil {
			return err
		}
		fmt.Printf("%s: %s\n", "Policy", policy.Name)
		fmt.Printf("%s: %s\n", "Action", policy.Desc())
	}
	query, err := resolveAddresses(cli, false)
	if err != nil {
		return nil
	}
	list, err := db.ListAcl(context.Background(), policyName, query)
	if err != nil {
		return err
	}

	if showHeader {
		fmt.Println("Addresses:")
	}
	for _, address := range list {
		fmt.Println(address.Hex())
	}
	return nil
}

func describePolicies(cli *cli.Context, showHeader bool) error {
	_, db, err := configAndStorage(cli)
	if err != nil {
		return err
	}
	list, err := db.DescribePolicies(context.Background())
	if err != nil {
		return err
	}

	if showHeader {
		fmt.Printf("%7s: %s\n", "Policy", "Action")
	}
	for _, p := range list {
		fmt.Printf("%7s: %s\n", p.Name, p.Desc())
	}

	return nil
}

func configAndStorage(cli *cli.Context) (*config.Config, *pgpoolstorage.PostgresPoolStorage, error) {
	c, err := config.Load(cli, false)
	if err != nil {
		return nil, nil, err
	}
	setupLog(c.Log)

	db, err := pgpoolstorage.NewPostgresPoolStorage(c.Pool.DB)
	if err != nil {
		return nil, nil, err
	}
	return c, db, nil
}

func requirePolicyAndAddresses(cli *cli.Context) (pool.PolicyName, []common.Address, error) {
	policy, err := resolvePolicy(cli)
	if err != nil {
		return "", nil, err
	}
	addresses, err := resolveAddresses(cli, true)
	if err != nil {
		return "", nil, err
	}
	return policy, addresses, nil
}

func resolvePolicy(cli *cli.Context) (pool.PolicyName, error) {
	policy := cli.String(policyFlag.Name)
	if policy == "" {
		return "", nil
	}
	if !pool.IsPolicy(policy) {
		return "", fmt.Errorf("invalid policy name: %s", policy)
	}
	return pool.PolicyName(policy), nil
}

func resolveAddresses(cli *cli.Context, failIfEmpty bool) ([]common.Address, error) {
	var set = make(map[common.Address]struct{})
	if cli.IsSet("csv") {
		file := cli.String(csvFlag.Name)
		fd, err := os.Open(file)
		if err != nil {
			return nil, err
		}
		defer func(fd *os.File) {
			_ = fd.Close()
		}(fd)

		fileReader := csv.NewReader(fd)
		records, err := fileReader.ReadAll()

		if err != nil {
			return nil, err
		}
		for _, row := range records {
			for _, cell := range row {
				hex := strings.TrimSpace(cell)
				set[common.HexToAddress(hex)] = struct{}{}
			}
		}
	}

	for _, a := range cli.Args().Slice() {
		a = strings.TrimSpace(a)
		a = strings.Trim(a, ",|")
		if !strings.HasPrefix(a, "0x") {
			a = "0x" + a
		}
		set[common.HexToAddress(a)] = struct{}{}
	}
	var ret []common.Address
	for a := range set {
		ret = append(ret, a)
	}
	if failIfEmpty && len(ret) == 0 {
		return nil, errors.New("no addresses given")
	}
	return ret, nil
}
