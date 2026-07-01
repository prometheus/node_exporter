# AGENTS.md

## Project Shape
- Main binary entrypoint is `node_exporter.go`; collectors live under `collector/` and self-register from `init()` via `registerCollector` in `collector/collector.go`.
- Collector availability is controlled by Go build tags and OS-specific filenames; check `//go:build` lines before assuming a collector exists on the current OS.

## Rules
- Use `kingpin.Flag()`, never `kingpin.Flag().Envar()`: the project has a "one-way-to-configure" policy of CLI flags only. Runtime env vars (`GOMAXPROCS`) and deprecated collectors are the only historical exceptions; maintainers have rejected env var flags for new code.
- Collector code must not shell out: `.golangci.yml` denies `os/exec` outside tests, and `CONTRIBUTING.md` says collectors may read `/proc`, `/sys`, use syscalls, or local sockets only.
- Keep metrics in the `node` namespace and update README collector tables when adding, removing, enabling, disabling, or changing collector flags.

## Project Map
- `node_exporter.go`: binary entrypoint, HTTP handler, collector filtering, and exporter bootstrap.
- `collector/`: collector implementations and tests; start here for metric behavior changes.
- `collector/collector.go`: collector registry, `NodeCollector`, scrape duration/success metrics, and `node` namespace.
- `collector/fixtures/`: test and e2e fixtures; `sys.ttar` and `udev.ttar` are source of truth for unpacked fixture dirs.
- `end-to-end-test.sh`: e2e golden generation/comparison and OS-aware collector flag filtering.
- `tools/`: helper binary used by e2e build-tag matching.
- `docs/node-mixin/`: separate Jsonnet alert/rule/dashboard generation workflow.

## Build And Test Commands
- Full CI-like local check is `make`; it runs vet, metric/rule checks, shared Prometheus checks, 32-bit tests when supported, and Linux e2e only.
- Focused unit tests: `go test ./collector -run TestName` or `go test ./... -run TestName`; use `make test` when tests need unpacked fixtures.
- `make test` unpacks `collector/fixtures/sys.ttar` and `collector/fixtures/udev.ttar`, then runs `go test -short $(test-flags) ./...`.
- `make build` uses `promu` and writes `./node_exporter`, which `make test-e2e` expects; root tests in `node_exporter_test.go` skip unless a binary exists at `$GOPATH/bin/node_exporter`.
- `make test-e2e` is Linux-only from the main Makefile path; it depends on `make build`, fixture unpacking, and `make tools`.
- `make checkmetrics` downloads/uses `promtool` and checks `collector/fixtures/e2e-output*.txt`; run it after metric output changes.
- Lint/format shortcuts: `make format`, `make lint`, `make lint-fix`, `make yamllint`, `make unused`.

## Fixtures And Generated Outputs
- Do not edit unpacked `collector/fixtures/sys/` or `collector/fixtures/udev/` as source of truth; update `collector/fixtures/sys.ttar` and `collector/fixtures/udev.ttar` with `make update_fixtures`.
- To refresh e2e metric golden files, run `make build tools` first, then `./end-to-end-test.sh -u`; it writes the appropriate `collector/fixtures/e2e-output*.txt`.
- `end-to-end-test.sh` filters collector flags by OS/build context using `tools/tools match`, so an ignored flag in e2e may mean the collector file is not compiled for that platform.

## Toolchain Notes
- CI uses Go 1.26 (`quay.io/prometheus/golang-builder:1.26-base` and `setup-go: 1.26.x`) even though `go.mod` currently says `go 1.25.0`; `.promu.yml` is Go 1.26 and `.promu-cgo.yml` is Go 1.25.
- `Makefile.common` is shared Prometheus infrastructure; comments say to open PRs against `prometheus/prometheus/Makefile.common` for generic changes.
- golangci-lint version comes from `make print-golangci-lint-version` (`GOLANGCI_LINT_VERSION` in `Makefile.common`), not from a locally installed default.
- Docker builds produce per-arch images named `$(DOCKER_REPO)/node-exporter-linux-<arch>:$(DOCKER_IMAGE_TAG)` and cover `Dockerfile` plus `Dockerfile.*` variants.

## Node Mixin
- Mixin checks are separate: CI runs `make promtool`, `make -C docs/node-mixin clean`, `make -C docs/node-mixin jb_install`, then `make -C docs/node-mixin` and expects `git diff --exit-code`.
- `docs/node-mixin` requires `jsonnet`, `jsonnetfmt`, and `jb`; generated files are `node_alerts.yaml`, `node_rules.yaml`, and `dashboards_out/`.
