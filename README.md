# mgraph

Graph manager

Notes
* used by md5 (quick solution) for hashing customer ids

### TODO list

- [ ] switch to `go mod`
- [ ] switch to `pgx`?
- [ ] add linter
- [ ] api
  - [ ] REST
  - [ ] gRPC
  - [ ] add integration tests
  - [ ] notifications - move_node, remove_node, add_node
- [ ] docs, translate all into en
- [ ] wrap in docker
- [ ] store inmemory (?), materialized paths or?
- [x] store postgres, materialized paths
  - [ ] more checks in pg functions
- [ ] store aws dynamodb, materialized paths (supports prefix search but not true transactions?)
- [ ] id the node sets the app?
  - [ ] id as string or generated uuidv4

## Tips

To change the scheme
```sql
SET search_path = mgraph, public;
```
