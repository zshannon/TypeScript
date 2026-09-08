# TypeScript compiler packages for Go

This module is generated from every package under `tsc/internal` in the
[TypeScript Go repository](https://github.com/zshannon/TypeScript). It publishes
the whole compiler package surface with normal Go visibility: exported identifiers
remain exported, while unexported identifiers remain private to their packages.

For example:

```go
import (
	"github.com/zshannon/TypeScript/public/v7/ast"
	"github.com/zshannon/TypeScript/public/v7/parser"
)
```

Nested source directories named `internal` are generated as `internals` so
packages such as `vfs/internal` can be imported by external modules.

## Compatibility

This broad API follows upstream compiler implementation packages and may change at
any time. Go module versions align exactly with upstream TypeScript versions. For
example, `github.com/zshannon/TypeScript/public/v7@v7.0.2` corresponds to TypeScript 7.0.2.

The current generated source is `7.1.0-dev`, a development version; it is not the stable 7.0.2 release.

## Regeneration

From the repository root:

```sh
go run ./scripts/export-go.go
go run ./scripts/export-go.go --check
go run ./scripts/export-go.go --version
```

The generator copies compiler packages and their required source assets. It does not
copy the repository-level `tsc/testdata` corpus; upstream test harness packages that
load those fixtures at runtime require a full TypeScript repository checkout.
