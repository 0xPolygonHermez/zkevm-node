# Component: Prover

NOTE: The Prover is not considered part of the XGON Node and all issues and suggestions should be sent to the [Prover repo](https://github.com/okx/xgon-prover/).

## XGON Prover:

The XGON Prover image hosts different components, *Merkle Tree*, *Executor* and finally the actual *Prover*.

## Hard dependencies:

- [Aggregator](./aggregator.md)

## Running:

The preferred way to run the XGON Prover component is via Docker and Docker Compose.

```bash
docker pull hermeznetwork/xgon-prover
```

To orchestrate multiple deployments of the different XGON Node components, a `docker-compose.yaml` file for Docker Compose can be used:

```yaml
  xgon-prover:
    container_name: xgon-prover
    image: xgon-prover
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