# Generating an Account Keystore file:

```bash
docker run --rm hermeznetwork/zkevm-node:latest sh -c "/app/zkevm-node encryptKey --pk=[your private key] --pw=[password to encrypt file] --output=./keystore; cat ./keystore/*" > account.keystore
```

**NOTE**:

- Replace `[your private key]` with your Ethereum L1 account private key
- Replace `[password to encrypt file]` with a password used for file encryption. This password must be passed to the Node later on via env variable (`ZKEVM_NODE_ETHERMAN_PRIVATEKEYPASSWORD`)
