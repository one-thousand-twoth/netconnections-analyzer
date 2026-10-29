# Network connection diagrams

**Status: planning.** This repository does not yet contain the Go executable. The planned `onethou-ncviz` CLI will turn a semicolon-delimited CSV of **declared connections** into a directed Graphviz `.dot` file. It will not discover network traffic or render images. Graphviz is a development/test dependency, not a production dependency.

The [MVP plan](docs/mvp.md) summarizes the input contract, diagram behavior, CLI, tests, and build order. The [completed decision map](https://github.com/one-thousand-twoth/netconnections-analyzer/issues/1) links the authoritative decision tickets; the [implementation handoff](https://github.com/one-thousand-twoth/netconnections-analyzer/issues/5) defines acceptance criteria.

## Planned input and output

```csv
src;dest;src_meta;dest_meta;edge_meta;ports;protocol
web-01;db-01;frontend;database;app traffic;5432;TCP
web-01;cache-01;;cache;session lookup;6379;TCP
```

Only `src` and `dest` are required. A valid row becomes a directed edge; repeated rows remain distinct edges. The default output for `connections.csv` will be `connections.dot`. Column names may be mapped explicitly with a JSON file.

**Planned command (not available yet):**

```sh
onethou-ncviz --input connections.csv
onethou-ncviz --input connections.csv --output diagram.dot --mapping columns.json --strict
```

Default mode warns and continues past rows with missing endpoints or conflicting host metadata. Strict mode fails on those errors without writing DOT. Malformed CSV and invalid mappings fail in either mode. See the [MVP plan](docs/mvp.md) for the full contract.

## Planned development workflow

Project commands will be mise tasks: `mise test` runs unit and Graphviz validation tests, and `mise test:graphviz` runs the Graphviz checks. Linux GitHub Actions will run the checks and cross-compile Windows binaries. Neither the tasks nor the CI workflow have been created yet.
