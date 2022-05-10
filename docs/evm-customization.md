# EVM Customization

## Modifications

* ChainID Opcode should return default chainID in all cases
* Sstore key and values must be passed in clear (not hashed) to the MT that will later hash then using poseidon


## Differences between EVM & zkEVM

| opcode       | EVM                                                         | ZKEVM                                                                                                                                                   |
| ------------ | ----------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------- |
| BLOCKHASH    | Get the hash of one of the 256 most recent complete blocks. | Get the state-root of a completed batch of transactions. If batch number given is further than the current batch - 1, the blockhash returned will be 0. |
| SELFDESTRUCT | Halt execution and register account for later deletion.     | Clears ONLY the account's bytecode from the trie, transfers remaining balance to the beneficiary (input 1) and halts the execution.                     |
| DIFFICULTY   | Get the block’s difficulty                                  | Always returns 0                                                                                                                                        |
