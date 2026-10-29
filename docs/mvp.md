# DOT-only MVP plan

This is a readable summary of **decisions**, not a description of implemented functionality. The [decision map](https://github.com/one-thousand-twoth/netconnections-analyzer/issues/1) is the index; its linked tickets hold the authoritative detail. Work begins from the [implementation handoff](https://github.com/one-thousand-twoth/netconnections-analyzer/issues/5). Use the terms in [`GLOSSARY.md`](../GLOSSARY.md).

## Goal and scope

Build a local, scriptable Go executable, `onethou-ncviz`, that reads assertions about connections between hosts from CSV and writes deterministic Graphviz DOT. It represents the input; it does not verify network connectivity. The MVP produces **only `.dot` files**. Production requires no Graphviz installation; development tests use Graphviz to validate the output. Image export, automatic column-name guessing, configurable graph styling, stdin/stdout streaming, and service APIs are out of scope. [Output decision](https://github.com/one-thousand-twoth/netconnections-analyzer/issues/7)

## Input contract

Default header (semicolon-delimited CSV, with a required header row):

```csv
src;dest;src_meta;dest_meta;edge_meta;ports;protocol
web-01;db-01;frontend;database;app traffic;5432;TCP
```

- `src` and `dest` are required, nonempty values for a valid connection. The five other fields are optional, plain text; ports and protocol are not interpreted as network configuration. Optional columns may be absent entirely.
- CSV headers match case-insensitively after trimming outer whitespace; host names and metadata values retain their case. Trim outer whitespace from cell values. Ambiguous matching headers are fatal. Ignore unmapped extra columns.
- A JSON mapping file may override any of the seven default header names. Its keys use the exact canonical field names, for example `{"src":"FromHost","dest":"ToHost"}`; CSV header matching for its values is still case-insensitive. Invalid JSON, unknown or duplicate keys, empty mapped names, and mappings to absent headers are fatal, including for optional fields. There is no automatic alias detection.
- Host metadata comes from `src_meta` and `dest_meta`. The first nonempty value encountered for a host on a valid row wins, regardless of whether that host appears as source or destination. Later conflicting values are handled by the selected error mode.
- An entirely blank row is skipped silently. In default mode, a nonblank row with an empty source or destination warns with its row number and is skipped; a conflicting host metadata value warns and leaves the first value in place. Warnings go to stderr and do not make the run fail. A valid-header file may yield an empty graph, even if every nonblank row was skipped.
- In strict mode, missing endpoints or conflicting host metadata are fatal and produce no DOT. Unreadable input, malformed CSV structure, and missing or ambiguous required mapped headers are fatal in **both** modes. [Input decision](https://github.com/one-thousand-twoth/netconnections-analyzer/issues/4)

## DOT semantics

- Use a directed, left-to-right `digraph` (`rankdir=LR`), **not** DOT's `strict digraph`. Every valid CSV row becomes one edge, including exact duplicates and self-loops. Host identity is case-sensitive.
- Nodes appear in first-seen order (source then destination for each row); edges appear in CSV row order. Identical input and configuration produce byte-identical DOT.
- A node label contains the host name, then `meta: …` on a new line if host metadata exists. An edge label contains only nonempty lines, ordered `ports: …`, `protocol: …`, `meta: …`. Keep an edge even if it has no label.
- Generate opaque node IDs. Encode host names and metadata as literal label text, including quotes, backslashes, newlines, `\N`-looking sequences and markup-looking strings. Do not turn CSV values into DOT syntax, arbitrary attributes, links, or images. [Graph decision](https://github.com/one-thousand-twoth/netconnections-analyzer/issues/3)

## Planned CLI

```text
onethou-ncviz --input connections.csv [--output diagram.dot] [--mapping columns.json] [--strict]
```

`--input` is required. Without `--output`, replace the input's last extension with `.dot`, or append `.dot` if it has no extension. An explicit output must end in `.dot`. A successful run replaces an existing DOT file. A failed run leaves any existing output unchanged and creates no partial new file. Reject the same normalized input/output path, without resolving symlink or hard-link aliases. Fatal errors return nonzero; warnings alone return zero in default mode. The intended targets are Linux and Windows, but this plan only validates Windows by cross-compilation, not by a runtime test. [CLI decision](https://github.com/one-thousand-twoth/netconnections-analyzer/issues/2)

## Modules and build order

Use module path `github.com/one-thousand-twoth/netconnections-analyzer`. Plan two importable packages outside `internal/`: a connection-data module taking an `io.Reader`, mapping and mode and returning connection/host data, warnings or errors; and a DOT module taking that data and an `io.Writer` to produce deterministic DOT. Document their interface invariants and error modes before implementation. Keep JSON file loading, CLI flags and output-file publishing private to the executable. Exact Go signatures may evolve during greenfield work before the first release.

1. Write CSV, mapping JSON, and expected DOT fixtures before parser implementation.
2. Initialize the Go module, executable skeleton and public module interfaces. Add pinned mise tooling and enforce formatting, vet, tests, and lint checks before feature code.
3. Implement CSV parsing, explicit mapping, default/strict behavior and host metadata precedence with unit tests.
4. Implement the graph/DOT writer with literal-text escaping, golden files, and repeat-run determinism tests.
5. Wire CLI paths, error reporting, successful output replacement and failure preservation; test them end to end.
6. Validate generated DOT with Graphviz in local and Linux CI tests; cross-compile Windows amd64.
7. After the MVP passes acceptance, manually merge the release-please PR and publish Linux amd64 and Windows amd64 binaries.

These are **planned stages**, not completed work. [Full handoff and exit criteria](https://github.com/one-thousand-twoth/netconnections-analyzer/issues/5)

## Planned development and release checks

- Pin Go, `golangci-lint`, and `aqua:gitlab.com/graphviz/graphviz` in mise; validate that a clean Linux `mise install` supplies `dot`. All documented developer commands and CI checks use mise tasks, not direct `go test` commands.
- `mise test` runs unit and Graphviz integration tests; `mise test:graphviz` must fail clearly if `dot` is missing, not silently skip. Test Graphviz parsing of generated DOT, not pixel-identical rendered images. Fixtures cover quoted delimiters, custom and ambiguous headers, blank/invalid rows, both error modes, conflicting metadata, duplicate edges, self-loops and DOT-special text.
- Linux GitHub Actions enforces mise tasks for format, vet, tests and `golangci-lint` config verification and run. Start with `linters.default: none` and **only `staticcheck` enabled**. Cross-compile Windows amd64; do not claim Windows runtime validation.
- Check PR titles and commits for Conventional Commits. Configure release-please to open a standard release PR; a human merges it only after MVP acceptance. The release workflow then attaches Linux amd64 and Windows amd64 binaries. [Acceptance and release decision](https://github.com/one-thousand-twoth/netconnections-analyzer/issues/5)
