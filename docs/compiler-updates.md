---
generated-by: GPT-6 Astra Ultra
---
# Compiler updates

This repository is a real public GitHub fork of `microsoft/TypeScript`. Its main page compares Git ancestry with that upstream. Compiler source updates use merge commits, preserving that ancestry even when we intentionally exclude upstream workflow changes.

## Source update PRs

The weekly/manual source-sync workflow fetches Microsoft `main`, merges our current `main` into a stable `automation/typescript-upstream` branch, merges upstream, regenerates `public/`, and opens or refreshes a review PR. Maintainer fixes on that branch are preserved. Pushes never force-overwrite a concurrent human update. Source conflicts stop the run for manual resolution.

Only our workflow files belong in `.github/workflows`. Each source merge preserves the exact workflow directory from our branch after incorporating our `main`. Upstream workflow additions, edits, and deletions are excluded; conflicts exclusively in that directory are resolved using our files. GitHub's generic Sync fork button does not implement this file policy; use the source-update PR workflow.

Review and merge source-update PRs with **Create a merge commit**, preserving their upstream merge parent. Nothing automatically merges PRs.

## Configure the update bot

1. Create a GitHub App owned by the repository owner with repository Contents, Pull requests, and Workflows read/write permissions. Install it only on `zshannon/TypeScript`.
2. Set repository Actions variable `UPSTREAM_SYNC_APP_ID` to its App ID and secret `UPSTREAM_SYNC_APP_PRIVATE_KEY` to its full PEM private key.
3. Merge the compiler validation workflows before enabling Actions so only our workflows are active in the default branch.
4. Manually run the source-sync workflow and confirm that its PR starts the compiler checks. An App token is used so normal PR workflows trigger; the default Actions token would suppress those downstream events.

The App installation and settings are separate repository setup steps. Preparing or pushing this branch does not establish them.

## Consume updates

Committed `public/` source is consumable by commit as soon as its branch is public. Consumers use normal `go get github.com/zshannon/TypeScript/public@<commit>`; Go records a pseudo-version for that commit. There is no custom tagging workflow or separate release numbering scheme. After review and merge, the server can update its pinned compiler dependency to that commit.

## Server migration

The old `zshannon/typescript-go` repository retains the server, native React compiler work in PR #15, and its existing deployment history. The next phase moves server code to its own module and replaces its direct compiler-internal imports with `github.com/zshannon/TypeScript/public/...`. Its esbuild and Oxc updates and Docker publishing belong there. This compiler branch contains no server or Docker deployment workflow.
