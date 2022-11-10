# Batchsender

## Overview

This script allows to send a specified number of (empty) batch transactions to
L1.  Optionally it can wait for the batch to be verified.  The script interacts
with L1 only. Basically it acts like a sequencer but without building a real
rollup.  It can be useful to test the Synchronizer and the Aggregator.

## Usage

The script can be installed running `go install` from this folder.

## Examples

- Send 1 batch:

```sh
$ batchsender 
```

- Send 1 batch and wait for its proof:

```sh
$ batchsender -w send
```

- Send 42 batches:

```sh
$ batchsender send 42
```

- Send 42 batches and wait for the proofs:

```sh
$ batchsender -w send 42
```

- Send 42 batches with verbose logs and wait for the proofs:

```sh
$ batchsender -v -w send 42
```
