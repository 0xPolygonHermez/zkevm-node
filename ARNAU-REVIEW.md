# Sequencer review

## Overall

- We need to do a deep review of L2 reorg handling

## Needs fix

- `func (f *finalizer) SortForcedBatches(fb []state.ForcedBatch) []state.ForcedBatch` should filter repeated items
- In `func (f *finalizer) handleSuccessfulTxProcessResp(ctx context.Context, tx *TxTracker, result *state.ProcessBatchResponse) error`, what if `len(result.Responses) == 0`? This should mean that the tx has not been executed isn't it? In fact it would be nice to add a sanity check: if stateRoot has not changed, it means that the tx has not been executed, therefore we need to reject it. [WIP] I got the feeling that this same logic is going to be used when porcessing empty/forced batch. Maybe we should treat this cases with different logic, as special cases, to avoid problems... After going further in the review: yes this happens, we should have a function to porcess empty batches, that should be super simple and dont care about tx responses and so
- [MAYBE] `func (w *Worker) MoveTxToNotReady(txHash common.Hash, from common.Address, actualNonce *uint64, actualBalance *big.Int)` could be called with both nonce and balance being nil. What to do in this case?
- It's a bit ugly to not check the error [here](sequencer/finalizer.go#L205). Even more if for whatever reason the tx triggers an unexpected error from the executor we would get in an endless loop until there is a tx with better efficiency on the worker. We need to review how to handle this kind of edge case
- this can be dangerous, it's possible that there are scenarios in wich we need to close a batch with same state root as previous one? Need to double check and think fuerther... Example: we need to close a batch with no txs and GER 0x0...0 as previous one, then the batch will have same state root. Apart from that the batch should have been already re-executed at this point (if it didn't had txs)
- `// TODO: design error handling for reprocessing`: fore an L2 reorg by dropping this batch and the following ones **if not sent yet to L1! we need to add something to mark the batch as ready to send to L1 when the sanity check is completed**
- `lastBatchNumberInState, stateRoot = f.processForcedBatch(lastBatchNumberInState, stateRoot, forcedBatch)` this must return an error, and clear only the forced batches that have actually been porcessed from the list of forced batches
- TODO: review `func (d *dbManager) ProcessForcedBatch(forcedBatchNum uint64, request state.ProcessRequest) (*state.ProcessBatchResponse, error)`. So far it looks like that the closing signal manager injects unnecessary data into the finalizer (raw txs data), because in the end we get this data again from the DB. In order to be safe what we could do is simplify all of this and the forced batch signal is just a boolean. If activated, we get all the pending forced batches (that have been pending for more than X blocks / that were added in txs that are finalized on L1) and process them. This is more simple, more safe, and probably can be implemented with even less queries than what we have right now

```go
// closeBatch closes the current batch in the state
func (f *finalizer) closeBatch(ctx context.Context) error {
	// We need to process the batch to update the state root before closing the batch
	if f.batch.initialStateRoot == f.batch.stateRoot {
		err := f.processTransaction(ctx, nil)
		if err != nil {
			return err
		}
	}
```

- resetting the deathline could be dangerous if we have an error and return before processing all the pending forced batches

```go
func (f *finalizer) processForcedBatches(ctx context.Context, lastBatchNumberInState uint64, stateRoot common.Hash) (uint64, common.Hash, error) {
	f.nextForcedBatchesMux.Lock()
	defer f.nextForcedBatchesMux.Unlock()
	f.nextForcedBatchDeadline = 0
```

## Need to improve

- The timer for executor usage is started before a lock. It should be started / stopped right before / after calling `result, err := f.executor.ProcessBatch(ctx, f.processRequest)`. In fact, this also includes a bunch of struct transformation and other logic that happens in the `state`. So it's not very accurate to meassure this way
- We need to understand why we do changes `finalizer -> worker` for instance "deleting tx bcs it uses all the counters in an empty batch"
- Missleading function:

```go
// checkRemainingResources checks if the transaction uses less resources than the remaining ones in the batch.
func (f *finalizer) checkRemainingResources(result *state.ProcessBatchResponse, tx *TxTracker) error
```
Apart of "checks if the transaction uses less resources than the remaining ones in the batch" it also "sub the used resources to `f.batch.remainingResources` if they're not underflown" AND "update the worker pool with used resources in case of underflow"

## Future improvements

- Change how the sequence sender works to be consistent with other workers (super low priority). Alternatively, just go ahead and implement it [this way](https://github.com/0xPolygonHermez/zkevm-node/issues/1631)
- Differentiate worker time when it's blocking vs when it's doing "background" job
- it's a bit crazy that we've to procees through the dbManager xD `response, err := f.dbManager.ProcessForcedBatch(forcedBatch.ForcedBatchNumber, processRequest)`

## Low priority / cosmetic improvements

- I've noticed that there is a `sequencer/mock` directory, but we have a bunch of `sequencer/mock_*.go` files. Could we move them to `sequencer/mock`?
- `func (f *finalizer) Start(ctx context.Context, batch *WipBatch, processingReq *state.ProcessRequest)` Probably shouldn't get batch & processingRequest as input param and always request it from DB?
- In `func (f *finalizer) processTransaction(ctx context.Context, tx *TxTracker) error`, there is two `if tx != nil {` checks. Maybe we should check only once? Actually is there any legit case where we should "processTransaction" a nil tx?
- `func (f *finalizer) handleSuccessfulTxProcessResp(ctx context.Context, tx *TxTracker, result *state.ProcessBatchResponse) error` is called before knowing if the tx is successful. I would either rename the func or check if the tx is succesful before calling the func