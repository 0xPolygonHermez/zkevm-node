# Component: Prover

## ZKEVM Prover:

The ZKEVM Prover image hosts different components, *Merkle Tree*, *Executor* and finally the actual *Prover*.

## Running:

The preferred way to run the ZKEVM Prover component is via Docker and Docker Compose.

```bash
docker pull hermeznetwork/zkevm-prover
```

To orchestrate multiple deployments of the different ZKEVM Node components, a `docker-compose.yaml` file for Docker Compose can be used:

```yaml
  zkevm-prover:
    container_name: zkevm-prover
    image: zkevm-prover
    volumes:
      - ./prover-config.json:/usr/src/app/config.json
    command: >
      zkProver -c /usr/src/app/config.json
```

The `prover-config.json` file contents will vary depending on your use case, the main document explains different values to be changed to achieve different behaviors.

### Ports:

- `50051`: Prover
- `50061`: Merkle Tree
- `50071`: Executor