---
generated-by: GPT-6 Astra Ultra
---
# Compiler packages

The module is `github.com/zshannon/TypeScript/public`. Package paths include `/compiler`, `/checker`, `/ast`, `/parser`, `/vfs`, and all other packages under upstream `tsc/internal`.

Generation copies whole package implementations and rewrites references between them. Standard Go visibility still applies to declarations inside a package. Nested directories named `internal` become `internals` in the generated tree, so even those packages can be imported outside the compiler module. For example, upstream `tsc/internal/vfs/internal` becomes `public/vfs/internals`.

The source tree remains under Microsoft's original module path and directory layout. It includes our compiler tracing and lifecycle fixes. Edit that source and regenerate; never hand-edit `public/`. Keeping the two trees separate makes incoming upstream changes ordinary source merges while giving consumers a conventional module path owned by this fork.

## Consume from another repository

Use Go 1.26 or newer. Select the exact commit reviewed in the compiler update PR, then run `go get github.com/zshannon/TypeScript/public@<commit>`. Go records the selected commit as a standard pseudo-version in `go.mod`; no separate tag or release process is needed.

Import packages directly, for example `github.com/zshannon/TypeScript/public/compiler`. No Microsoft module replacement or matching internal import prefix is needed. The executable example under `scripts/testdata/public-consumer` parses and typechecks TypeScript with embedded standard libraries, including a deliberate type error assertion.

Pin dependency versions because these are upstream compiler packages, not a separately stabilized wrapper API. Go downloads the module source directly from the public Git repository or its module proxy; there is no separate package upload service.

## Generate and validate

From the repository root:

```sh
go run ./scripts/export-go.go
go run ./scripts/export-go.go --check
go test ./scripts/export-go.go ./scripts/export-go_test.go
bash scripts/test-public-module.sh
```

The generator preserves platform build constraints, assembly, licenses, and embedded runtime assets. Its check mode detects missing, changed, and stale generated files without rewriting them. The consumer test stages the module outside this checkout, uses an unrelated module identity and `GOWORK=off`, and exercises real parsing and typechecking.

Upstream test harness packages are included with the rest of the compiler packages. Their repository-scale fixtures remain in upstream `tsc/testdata`; use the source checkout for upstream's complete compiler test suite. Runtime embedded standard libraries and diagnostics are included in the public module.
