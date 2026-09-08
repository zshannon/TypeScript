---
generated-by: GPT-6 Astra Ultra
---
# Compiler updates

This repository is a real public GitHub fork of `microsoft/TypeScript`. Its main page compares Git ancestry with that upstream. Compiler source updates use merge commits, preserving that ancestry even when we intentionally exclude upstream workflow changes.

## Source update PRs

The weekly/manual source-sync workflow fetches Microsoft `main`, merges our current `main` into a stable `automation/typescript-upstream` branch, merges upstream, regenerates `public/`, and opens or refreshes a review PR. Maintainer fixes on that branch are preserved. Pushes never force-overwrite a concurrent human update. Source conflicts stop the run for manual resolution.

Only our workflow files belong in `.github/workflows`. Each source merge preserves the exact workflow directory from our branch after incorporating our `main`. Upstream workflow additions, edits, and deletions are excluded; conflicts exclusively in that directory are resolved using our files. GitHub's generic Sync fork button does not implement this file policy; use the source-update PR workflow.

Review and merge source-update PRs with **Create a merge commit**, preserving their upstream merge parent. Nothing automatically merges PRs.

## Enable the automation

Merge the setup PR into `main`. The weekly workflow runs every Monday at 17:00 UTC and can also be run manually. In Settings → Actions → General, enable **Allow GitHub Actions to create and approve pull requests**. No GitHub App or additional secrets are required.

The bot uses the built-in repository token to maintain its upstream PR. It explicitly dispatches compiler validation for the update branch, so validation does not depend on automatic triggering from bot-created PR events. It never approves or merges its own PR.

## Releases from main

Every push to `main`, including each merged PR, runs the release workflow. It tests the compiler and external Go consumer before tagging the exact validated commit. The version comes from the compiler source: TypeScript `7.0.2` becomes Go `v7.0.2`; development source `7.1.0-dev` becomes `v7.1.0-dev`.

Stable versions are checked against the matching official Microsoft release tag. Prerelease versions are checked against the common upstream source ancestor. The public module uses the matching major suffix, such as `github.com/zshannon/TypeScript/public/v7`, and Go's standard submodule tag, such as `public/v7.0.2`.

If that TypeScript version already has a public tag, the run leaves it unchanged and skips publication. Go versions are immutable: multiple merges with the same TypeScript version do not produce additional releases. There is no independent fork version counter.

The normal flow is: review an upstream PR, merge it with a merge commit, and let the release workflow run. There is no separate release branch to maintain. Manual release dispatch remains available for a reviewed historical release snapshot.

## Server migration

The old `zshannon/typescript-go` repository retains the server, native React compiler work in PR #15, and its existing deployment history. The next phase moves server code to its own module and replaces its direct compiler-internal imports with `github.com/zshannon/TypeScript/public/v7/...`. Its esbuild and Oxc updates and Docker publishing belong there. This compiler branch contains no server or Docker deployment workflow.
