# Stable compiler updates

This fork publishes the Go compiler packages from stable Microsoft TypeScript releases only. The Go version matches TypeScript: upstream `v7.0.2` produces `public/v7.0.2` for `github.com/zshannon/TypeScript/public/v7`. There is no independent compiler version counter.

## Weekly updates

The scheduled **Sync stable TypeScript releases** workflow runs Renovate with the GitHub Releases datasource. Renovate tracks `.github/upstream-version`, ignores prereleases, and opens a draft PR when a newer stable release exists.

The same workflow fetches the exact release tag, merges its commit into the Renovate branch, preserves this fork's GitHub workflows, and regenerates `public/`. It then marks the PR ready and starts compiler CI. The update does not merge the current tip of upstream `main`.

Review and merge the PR using **Create a merge commit**, so the upstream release remains an ancestor of the fork. Do not squash or rebase upstream-update PRs. Renovate's automatic rebasing is disabled so it does not discard the source merge commits added to its branch. Conflicts leave the PR unmerged for review.

## Publishing

After a merge to `main`, **Release public Go module** validates the compiler version against Microsoft's exact stable tag, tests the compiler and public module, and publishes its matching immutable Go module tag. The workflow refuses dev, alpha, beta, and RC versions. Already published versions are left unchanged; fork-only changes need a new upstream version before another matching version can be published.

The historical `public/v7.1.0-dev` tag is not a supported release channel. It is retained as history and is not moved or relabeled as stable.

## Server

The server lives in the private `zshannon/tsgo-server` repository. Its dependency-update PRs consume the stable public module versions. The server's own SemVer, esbuild/Oxc updates, image publishing, and Fly/exe.dev deployments are independent of this compiler repository.

## Development branch

The same weekly workflow independently merges Microsoft’s `main` into `development`, preserving the fork’s workflows and regenerating its public packages. Merge conflicts fail the job without pushing a partial update. Development updates never publish a release.

Stable release PRs target our `main` and merge the exact upstream release tag, not the development branch. This excludes unreleased commits even when Microsoft’s `main` has advanced past the stable release.
