# Diff

This repo is a fork from [zkevm-node](https://github.com/0xPolygonHermez.zkevm-node). The purpose of the fork is to implement tha Validium consensus, enabling data availability to be posted outside of L1.

In order to document the code diff the [diff2html-cli](https://www.npmjs.com/package/diff2html-cli) tool is used. An html file is included in the repo [here](./diff.html). This file has been generated running the following command:

```bash
PATH_TO_ZKEVM_NODE_REPO="/change/this"
diff -ruN \
-I ".*github.com\/0x.*" \
-x "*mock*" -x ".git" \
-x ".github" \
-x ".gitignore" \
-x ".vscode" \
-x "ci" \
-x "environments" \
-x "*.md" \
-x "*.html" \
-x "*.html" \
-x "*.json" \
-x "*.toml" \
-x "*.abi" \
-x "*.bin" \
-x "*.pb.go" \
-x "smartcontracts" \
-x "go.sum" \
-x "mock*.go" \
-x "*venv*" \
-x "/dist/" \
-x "/test/e2e/keystore" \
-x "/test/vectors/src/**/*md" \
-x "/test/vectors/src/**/*js" \
-x "/test/vectors/src/**/*sol" \
-x "/test/vectors/src/**/*sh" \
-x "/test/vectors/src/package.json" \
-x "/test/contracts/bin/**/*.bin" \
-x "/test/contracts/bin/**/*.abi" \
-x "/tools/datastreamer/*.bin" \
-x "/test/datastreamer/*.db/*" \
-x "/test/*.bin" \
-x "/test/*.db/*" \
-x "**/.DS_Store" \
-x ".vscode" \
-x ".idea/" \
-x ".env" \
-x "out.dat" \
-x "cmd/__debug_bin" \
-x ".venv" \
-x "*metrics.txt" \
-x "coverage.out" \
-x "*datastream.db*" \
${PATH_TO_ZKEVM_NODE_REPO} . | \
diff2html -i stdin -s side -t "zkEVM node vs CDK validium node</br><h2>zkevm-node version: v0.5.0<h2/>" \
-F ./docs/diff/diff.html
```

Note that some files are excluded from the diff to omit changes that are not very relevant

## Policy aka Allow list

This is a new feature introduced for Validium (that could exist on rollup). `TODO: ` add more explanation / link to docs.

## Smart contracts

Currently using [this version of the smart contracts](https://github.com/0xPolygonHermez/zkevm-contracts/tree/feature/build-v4.0.0-rc.5-fork.7/contracts/v2/consensus/validium)

The main changes on Validium vs Rollup consensus smart contracts are:

- Validium is a superset of rollup, as it also implements sequencing as a rollup (disabled by default, but admin can enable this option)
- Validium implements a specific sequencing method that:
    - the function expects a new parameter: `bytes calldata dataAvailabilityMessage`
    - For each batch instead of sending a byte array of arbitrary size for data, it sends a hash, that must be the hash of what would have been the transactions on the rollup (if this is not the case the ZKP won't be verified later on)
    - It requires to call another contract, this contract needs to be set calling `setDataAvailabilityProtocol` before sequencing
    - A new hash is computed during the for loop that iterates over the batches, basically `accumulatedNonForcedTransactionsHash = keccak256(abi.encodePacked(accumulatedNonForcedTransactionsHash, currentBatch.transactionsHash));`
    - Finally the DA contract is called to verify that the DA has the data (note that the DA protocol is just an intrface, and this could be implemented in multiple ways. Currenty only Data Availability Committee (DAC) is implemented): `dataAvailabilityProtocol.verifyMessage(accumulatedNonForcedTransactionsHash, dataAvailabilityMessage);`

## Diff explained by package

### Cmd

- `newEtherman` now depends on the state, needed to instantiate the DA
- `LoadAuthFromKeyStore` returns the raw private key, used to authenticate messages for DAC
- policy CLI added to interact with allow list storage
- needed to change the order of how things are instantiated by adding a `tmpEthMan` as etherman is used to get the L2 chain ID, which is needed for the state, and at the same time the state is needed for for the Etherman. Note that the `tmpEthMan` is used exclusively to get the L2 chain ID
- `newEtherman` calls `newDataAvailability` as DA is a dependency of Etherman (etherman now pulls data from the DA)
- `newDataAvailability` used to instantiate the DA, in a modular way (switch/case of the supported backends)
- `createSequenceSender` instantiates DA as it's needed to post sequences to the DA layer

### Config

- removed hardcoded config options for testnet / mainnet / cardona. Those are meant to easilly config a netwrok of Polygon zkEVM Mainnet / Testnet Beta, but there's no such thing for Validium networks at the moment
- duration marshlling: used at some point when this type was used in the API of the DAC. This is likely not used anymore and could be removed

### Data Availability

This package is where most of the Validium logic happens, and it does not exist on `zkevm-node`, it has two main purposes:

- Retrieve data while synchronizing. This is done by:
    - Queryng the local state (in case the batch was already synced in the trusted state sync process)
    - Querying the trusted sequener
    - Querying the DA backend. In this case an interface is used to support arbitrary backends. Currently only DAC supported
- Post data to the backend: send data to the DA backend also using the backend intrface to be modular

### DB / Migrations

Migrations are named as an ascending number on `zkevm-node`. On validium, the prefix `validium-` is used when adding migrations specific to it, to avoid interferences whith the upstream migrations that get added.

Currently there is only one migration specific to Validium `db/migrations/pool/validium-001.sql` that is used for the policy feature

### Etherman

This is one of the most affected components in terms of code diff, as Etherman's purpose is to abstract the L1 interactions, and this interactions differ rollup vs validium.

- Add `da dataavailability.BatchDataProvider` to the main struct, to be able to do the DA query when synchronizing the virtual state (got a hash from L1, get the pre-image from that hash on the DA)
- The data availability protocol contract is instantiated in the `New` function to be able to implement functions realted to this contract. Note that:
    - The address of this contract is loaded from the mai Validium contract, so no extra config is needed
    - This is an interface of a contract, so no matter what actual DA backend implementation is used (currently only DAC)
- Delete code for `updateEtrogSequence`: this code is used to move from single LxLy to uLxLy. We expect to only do that for Hermez zkEVM, which is not a validium, hence I didn't bothered to fix this (it was failing is left as it is)
- `EstimateGasSequenceBatches`, `BuildSequenceBatchesTxData` and `sequenceBatches` now have an extra parameter `dataAvailabilityMessage []byte` to be consistent with the Validium call for `sequenceBatches` (`SequenceBatchesValidium`)
- `sequenceBatches` hashes the `BatchL2Data` to cofrom with the Validium way of sequencing batches on-chain. Up until this point the Etherman handles sequences as if they were from a rollup, so other components may remain less affected in terms of diff. But at this point is needed to do the "validium transformation" before actually interacting with L1
- `decodeSequences`: added a switch/case to support both rollup and validium decoding. In the Validium decoding, the DA is used to resolve the hash and get transactions data. By doing so, the rest of the code base is not aware of using validium or rollup (so the synchronizer can have a 0 code diff). The rollup decoding is not tested, and for sure it's not working. At very least the name on the switch case should be updated as explained on the `TODO:`. Note that the validium contract supports behaving as a rollup, but the node does not implement this yet
- Added getters for DA smart contract so other components can still use Etherman as an L1 abstraction
- On the test side of things:
    - added mocks for the DA
    - deployed the DAC on the simulated Etherman. This is something that could be reviewed to reflect the modular mindset, although the specific logic of the DAC is not tested. Note that the main Validium contract needs to have an instance of a contract that implements the interface, otherwise it will fail to `sequencevalidiumBatches` as it will call an address without code

### JSON RPC

Implements the allow list (policy) filtering

### Pool

Implements the allow list (policy) storage

### Sequencer

Remove unused functions from the interface. In general, the idea is try to reduce code diff at all cost, but in this case, the removed functions would have need to be updated (to include `dataAvailabilityMessage []byte`). This methods should be removed from the upstream as well

### Sequence Sender

That's one of the most relevant packages in terms of code diff, as the sequence sender is in charge of producing the virtual state. This means that it needs to do the sequence batches interaction with L1, which is different rollup vs validium.

Conceptually the main differences are:

- How the sequence is cut (how many batches are sent together): In rollup, calls to L1 are being done to simulate if the sequence will be accepted. If this is the case, the sequence sender will try to add another batch, until it starts failing (among other checks related to timeouts, and other logic that comes mostly through config). However, in Validium there is the need to interact with arbitrary Backends before interacting with L1 (for instance, in the case of DAC, the signatures of the committee members are needed to be able to successfuly call L1). This makes it not feasible to use L1 as a simulator in an iterative way, as it would spam the DA layer. Apart from that, rollup is more unpredictable since each batch can be very different in terms of calldata used, while in Validium this is fixed, as only a hash of the batch data is sent. For this reason Validium simplifies the trial / error logic for a simple check that will just try to group a fixed amount of batches `MaxBatchesForL1` (new config parameter). Other checks that don't need to interact with L1 are kept
- How the L1 tx is built: the L1 tx needs the `dataAvailabilityMessage`. This piece of data should be obtained from a DA backend. So the sequence sender now needs a `da dataAbilitier` interface to work. This interface will allow to collect the message after posting the data from the sequence. How the message looks like or how the data is posted to the backend is up to each protocol implementation (currently only DAC supported)

### State

Add `GetBatchL2DataByNumber` method to be able to retrieve the batch data without other parts of the batch

### Test

The docker image with the L1 with pre-deployed contracts is different, as it uses a Validium consensus contract. In consequence config files with different addresses and deploy block are used.

Add a test specific for DAC, this requires new containers configured in a way to force using the DA backend (instead of the trsuted sequencer) for synchronization. Also adds a DB for the DAC nodes. In order to implement this, changes are done in:

- Makefile
- docker-compose
- The new test file
- Operations manager

### Upstream Version

v0.5.0