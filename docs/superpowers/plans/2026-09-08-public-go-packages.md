---
generated-by: GPT-6 Astra Ultra
---
# Public compiler packages implementation plan

**Goal:** Make every TypeScript compiler package importable from ordinary Go modules.
**Architecture:** Generate `public/` from upstream `tsc/internal` with consistent import rewriting, retaining the original compiler tree and fork ancestry. Keep server/deployment outside this compiler fork.
**Tech stack:** Go 1.26, Go modules, GitHub Actions, existing Bash upstream-sync tests.
**Spec:** `docs/superpowers/specs/2026-09-08-public-go-packages.md`

## Constraints

1. Module `github.com/zshannon/TypeScript/public`; all compiler packages, no symbol allowlist.
2. Preserve assets, platform/build constraints, and compiler patches; rewrite nested `internal` segments to `internals` only in generated output.
3. Deterministic committed output with stale-file deletion and check mode; no source-tree relocation.
4. Genuine public fork and merge ancestry; only our workflows; no auto merge or deployment.

## Task 1: Generate the public module

Files: `scripts/export-go.go`, `scripts/export-go_test.go`, generated `public/**`.
Interface: `go run ./scripts/export-go.go` generates; `go run ./scripts/export-go.go --check` verifies without modifying.
- [ ] Write fixture tests covering cross-package imports, nested internal directories, embedded/binary assets, build constraints, stale output removal, idempotence, and check mode.
- [ ] Implement the deterministic whole-package projection; generate go.mod/go.sum/license and consumption documentation.
- [ ] Run `go test ./scripts/export-go.go ./scripts/export-go_test.go`, check mode, and `GOWORK=off go build ./...` inside public.

## Task 2: Isolate compiler CI and automate upstream updates

Files: `.github/workflows/**`, `.github/scripts/**`, `scripts/test-public-module.sh`, external smoke example under `scripts/testdata/public-consumer`, README and compiler publication docs.
Consumes: Task 1 command and module path.
- [ ] Keep only compiler package validation and upstream sync workflows.
- [ ] Verify consumption in an unrelated temporary module, with real parse/typecheck assertions and standard libraries.
- [ ] Reuse tested upstream workflow-ownership merge scripts; regenerate/commit public output before update PR push.
- [ ] Consume the public module through ordinary Go commit pinning; create no release workflow or custom tags.
- [ ] Run script/YAML checks, existing sync integration cases, source compiler regressions, module consumer tests, and representative cross-platform builds.

## Task 3: Review and publish

- [ ] Review generator, external-consumer proof, workflow ownership/regeneration, and release behavior together; fix material findings.
- [ ] Commit and save the branch to the original repository without changing its checkout.
- [ ] Push to the real public fork and open a draft PR; verify repository parent, visibility, and remote commit SHA. Do not merge or activate deployment.
