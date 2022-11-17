# Sequencer v1

## Flow

1. jRPC executes to get *status


- *status

```go
type TxStatus struct {
    ExecutedAtRoot common.Hash
    WasExecutedSuccessfuly bool
}
```