---
generated-by: GPT-6 Astra Ultra
---
# Public compiler Go packages

The user approved separating the server from the TypeScript fork and explicitly requested broad publication of compiler internals, rather than a curated adapter. This first phase publishes the compiler packages; server extraction follows separately.

Use the real public `zshannon/TypeScript` fork of `microsoft/TypeScript`, retaining upstream ancestry and GitHub's native ahead/behind comparison. Base this work on current fork main `d0ac85d8eab09e06376390c9352d07c8a8c00aa8`. Carry forward the eleven compiler tracing/lifecycle patch files from the tested port; do not import the server, Oxc, esbuild, or Docker deployment into this compiler-only branch.

Generate a Go module at `public/`, module path `github.com/zshannon/TypeScript/public/v7` (major derived from the compiler version), from all packages under `tsc/internal`. Expose packages such as `public/ast`, `public/compiler`, and `public/checker`, preserving their exported declarations and implementation. Rewrite references among compiler packages consistently. Rename nested `internal` path segments to `internals` in the generated copy so every package can be imported externally. Keep upstream source paths unchanged. Preserve build constraints, assembly, licenses, and embedded assets. No hand-maintained symbol allowlist, API wrapper layer, or import-path impersonation.

Generation must be deterministic, remove stale generated files, include a read-only check mode, and remain repeatable after commits. Commit generated packages so standard Go module resolution works from a commit. Test generation edge cases and prove consumption from a temporary module with an unrelated module path and `GOWORK=off`, including real parsing/typechecking with embedded standard libraries. Build every public package; verify relevant Linux and Darwin platform builds.

The fork retains only our compiler publication/test and source-sync workflows. Source sync preserves our workflow directory and regenerates the public module in the same update PR. No automated merging. The user requires the Go package version to match the upstream TypeScript version exactly. Derive the version and module major from compiler source. Use standard Go submodule tags such as `public/v7.0.2` for TypeScript `7.0.2`, with no independent fork version series. Validate the matching official release source before publishing; never label development main as an older stable release. GitHub hosts Go module source directly, so no registry upload is needed. Publishing a reviewed branch and draft PR is authorized, including previously approved 1Password authentication. Do not merge the PR or change production deployment.
