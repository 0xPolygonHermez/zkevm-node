# Sequencer improvement

The sequencer is a key component of the node and the network as it gives a first level of finality which will be percieved by the users.
This document pretends to describe an architecture that helps improving the performance of it.

## Goals

> tl;dr: As per the current Ethereum behaviour (pre-EIP4844), we should try to hit **~45 kB/s throughput** while maximizing collected fees

- Maximize profit: sequencer main goal as an actor is to maximize profit. At this point in time MEV techniques are out of the equation. Therefore, the only way to maximize profit if by optimizing fee collection and reducing operational cost. In this iteration reducing operational cost (having as many transactions / batch, L1 tx gas price, ...) will not be taken into consideration deeply
- Throughput: how many transactions per second can we process? How many time a user needs to wait from tx sent to tx added into the state? Sequencer duty is to be as fast as possible to bring a great UX into the network. Since there is a hard limit imopsed by L1 (Ethereum), the sequencer should try to reach this limit.

Ethereum imposed throughput limit:


```
Max bytes per block = X
Gas consumption = X * gas cost per byte + other gas consumption from the SC call
Max gas consumption = Block gas limit * percentage of the block the zkEVM shouldn't exceed
Max gas consumption = Gas consumption
Throughput limit = X / block time
```

Given the previous equation, with the following numbers:

- Gas cost per byte (call data) = K = 16 gas/Byte
- Other gas consumption from the SC call ~= Z ~= 250.000 gas
- Block gas limit = M = 30.000.000 gas
- Percentage of the block the zkEVM shouldn't exceed = P = 0.3
- Block time = T = 12s

The max throughput is

```
(M*P-Z)/(K*T) = Max Bytes/s
```

```
(30.000.000*0.3 gas -250.000 gas) / (16 gas/Byte * 12s) = 45kB/s
```

## Assumptions and perspective

- All transactions **must** be processed sequentialy before considering them final. Therfore in order to reach the goals described above, the design should aim at:
  - Having a great hit ratio (successfuly processed txs / total processed txs)
  - Minimizing time between txs are processed (as soon as one tx has been processed the next one starts)
  - Txs are sorted by max gas price
- Although there is a point in the system which can't be run in paralel there is no reason to have a system which is 100% sequential
- Processing transactions is a job done by the [executor](https://github.com/0xPolygonHermez/zkevm-prover) and it's performance is out of the scope of this document. An awesome and super fast performance will be assumed
- The system should be able to resume operations after any of it's components crashes. This is an important consideration whenever thinking wich pieces of data will be stored only in memory

## Architecture propposal

In order to follow the full concept, different parts of the system will be introduced separetly, and then will be put together at the end of the document.

Note that the different ideas are meant as abstract and don't suggest infrastructure or implementation details (unless the opposite is said). For instance, in the following diagrams shapes that usually represent DBs are used, but in this case they could be implemented using in memory structures, messaging systems (such as Apache Kafka) or actual DBs

### Finalizer

This is the component responsible for executing transactions sequentially and giving finality

![img](sequencer-finalizer.drawio.png)

The finalizer will consume data from 3 sources, after each interaction it will query this 3 sources in order:

- `Forced Batch Manager`: will indicate when a forced batch needs to be added in the state. When this happens the finalizer will close the open batch, process the forced batch(es) and finaly open a new batch
- `GER manager`: will indicate when a new GER or timestamp needs to be used. When this happens the finalizer will close the open batch, and open a new batch with the updated GER and timestamp
- `Ready Manager`: it's responsible for giving the next transaction to be executed, and should be the most frequently used source. Note that the delay between a request and a response for the next transaction to be processed and the hit ratio of the suggested transactions are a key part of the system.

In order to execute the transactions the `Super Executor` will be used. This is actually a normal executor, but it has this name to help understand that this specific instance needs to be the one performing the best, and should be used exclusively by the finalizer.

After (successfuly) executing transactions the finalizer will:

- Update the state with the new items (add to `StateDB`)
- Append the new `state mod logs`

On the other hand if the execution fails (this case is only considered when consuming data from the `Ready Manager` source), the failed transactions will be sent to the `Blocked Manager` (note that this component is ommited in the diagram, more on this later)

> QUESTION: is it really possible to update GER, timestamp and process transactions **without being aware of the batch number**

### State mod log

This component is a sorted list of logs that communicate changes on the state. The main purpose of it it's to help other components understanding if a given transaction can be processed without actualy executing it.

The content of this log should be something like this

```json
{
    "0xOldRoot0": {
        "nextRoot": "0xOldRoot1",
        "0xAddr0": {
            "nonce": 1337,
            "balance": "100000000000000000000000000",
            "storageModified": false
        },
        "0xAddrN": {
            "nonce": null,
            "balance": "0",
            "storageModified": true
        } 
    },
    "0xOldRootN": {
        "nextRoot": "0xOldRootN+1",
        "0xAddr0": {
            "nonce": 7331,
            "balance": "22222222222222",
            "storageModified": false
        },
        "0xAddrN": {
            "nonce": null,
            "balance": "0",
            "storageModified": true
        } 
    }
}
```

Note that:

- 