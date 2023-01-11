# Txsender

## Overview

This script allows to send a specified number of transactions to either L1 or
L2 (or both).  Optionally it can wait for the transactions to be verified.

## Usage

The script can be installed running `go install` from this folder.

## Examples

- Send 1 transaction on L2:

```sh
$ txsender 
```

- Send 1 transaction on L2 and wait for it to be validated:

```sh
$ txsender -w send
```

- Send 42 transactions on L1:

```sh
$ txsender -n l1 send 42
```

- Send 42 transactions both on L1 and L2 with verbose logs and wait for the validations:

```sh
$ txsender -v -w -n l1 -n l2 send 42
```
