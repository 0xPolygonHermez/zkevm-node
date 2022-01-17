package main

import (
	"fmt"
	"strings"

	"github.com/hermeznetwork/hermez-core/config"
	"github.com/hermeznetwork/hermez-core/db"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/state"
	"github.com/hermeznetwork/hermez-core/state/pgstatestorage"
	"github.com/hermeznetwork/hermez-core/state/tree"
	"github.com/iden3/go-iden3-crypto/poseidon"
	"github.com/jackc/pgx/v4"
	"github.com/urfave/cli/v2"
)

func registerSequencer(ctx *cli.Context) error {
	configFilePath := ctx.String(flagCfg)
	network := ctx.String(flagNetwork)
	url := ctx.Args().First()
	fmt.Print("*WARNING* Are you sure you want to register " +
		"the sequencer in the rollup using the domain <" + url + ">? [y/N]: ")
	var input string
	if _, err := fmt.Scanln(&input); err != nil {
		return err
	}
	input = strings.ToLower(input)
	if !(input == "y" || input == "yes") {
		return nil
	}

	c, err := config.Load(configFilePath, network)
	if err != nil {
		return err
	}

	setupLog(c.Log)

	runMigrations(c.Database)

	//Check if it is already registered
	etherman, err := newEtherman(*c)
	if err != nil {
		log.Fatal(err)
		return err
	}
	sqlDB, err := db.NewSQLDB(c.Database)
	if err != nil {
		log.Fatal(err)
		return err
	}
	store := tree.NewPostgresStore(sqlDB)
	mt := tree.NewMerkleTree(store, c.NetworkConfig.Arity, poseidon.Hash)
	tr := tree.NewStateTree(mt, []byte{})

	stateCfg := state.Config{
		DefaultChainID: c.NetworkConfig.L2DefaultChainID,
	}

	stateDb := pgstatestorage.NewPostgresStorage(sqlDB)
	st := state.NewState(stateCfg, stateDb, tr)

	_, err = st.GetSequencer(ctx.Context, etherman.GetAddress())
	if err == pgx.ErrNoRows { //If It doesn't exist, register the sequencer
		tx, err := etherman.RegisterSequencer(url)
		if err != nil {
			log.Error("uff no: ", err)
			return err
		}
		log.Info("Sequencer registered. Check this tx to see the status: ", tx.Hash())
		return nil
	} else if err != nil {
		return err
	}
	// If Sequencer exists in the db
	fmt.Print("*WARNING* Sequencer is already registered. Do you want to update " +
		"the sequencer url in the rollup usign the domain <" + url + ">? [y/N]: ")
	if _, err := fmt.Scanln(&input); err != nil {
		return err
	}
	input = strings.ToLower(input)
	if !(input == "y" || input == "yes") {
		return nil
	}

	tx, err := etherman.RegisterSequencer(url)
	if err != nil {
		return err
	}
	log.Info("Sequencer updated. Check this tx to see the status: ", tx.Hash())

	return nil
}
