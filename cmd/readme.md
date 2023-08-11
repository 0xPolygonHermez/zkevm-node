## Create and restore snapshots

### Create snapshots
```
go run ./cmd snapshot --cfg config/environments/local/local.node.config.toml --output ./
```

### Restore snapshots
```
go run ./cmd restore --cfg config/environments/local/local.node.config.toml -is zkevmpubliccorestatedb_1685614455_v0.1.0_undefined.sql.tar.gz -ih zkevmpublicstatedb_1685615051_v0.1.0_undefined.sql.tar.gz
```