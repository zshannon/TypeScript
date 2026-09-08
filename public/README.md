# TypeScript compiler packages for Go

This module is generated from every package under `tsc/internal` in the
[TypeScript Go repository](https://github.com/zshannon/TypeScript). It publishes
the whole compiler package surface with normal Go visibility: exported identifiers
remain exported, while unexported identifiers remain private to their packages.

For example:

```go
import (
	"github.com/zshannon/TypeScript/public/ast"
	"github.com/zshannon/TypeScript/public/parser"
)
```

Nested source directories named `internal` are generated as `internals` so
packages such as `vfs/internal` can be imported by external modules.

## Compatibility

This broad API follows upstream compiler implementation packages and may change at
any time. Consumers should pin a commit using its Go pseudo-version instead of
assuming semantic API stability. The module uses ordinary Go module resolution
from repository commits and has no separate release process.

## Regeneration

From the repository root:

```sh
go run ./scripts/export-go.go
go run ./scripts/export-go.go --check
```

The generator copies compiler packages and their required source assets. It does not
copy the repository-level `tsc/testdata` corpus; upstream test harness packages that
load those fixtures at runtime require a full TypeScript repository checkout.
