#!/usr/bin/env bash
set -euo pipefail
# Renovate owns release detection and PR creation. This step supplies the fork's
# merge commits and generated packages, which a version-file update cannot do.
trusted_merge=$(mktemp)
trap 'rm -f "$trusted_merge"' EXIT
cp .github/scripts/merge-typescript-upstream.sh "$trusted_merge"
git config user.name 'github-actions[bot]'
git config user.email '41898282+github-actions[bot]@users.noreply.github.com'
git remote add stable-upstream https://github.com/microsoft/TypeScript.git
gh pr list --repo "$GITHUB_REPOSITORY" --base main --state open \
  --json number,headRefName,isCrossRepository \
  --jq '.[] | select(.isCrossRepository == false) | select(.headRefName | startswith("renovate/typescript-")) | [.number,.headRefName] | @tsv' \
  > /tmp/typescript-update-prs
while IFS=$'\t' read -r number branch; do
  [[ -n "$number" ]] || continue
  git fetch origin "+refs/heads/main:refs/remotes/origin/main" "+refs/heads/$branch:refs/remotes/origin/$branch"
  git checkout -B "$branch" "origin/$branch"
  tag=$(cat .github/upstream-version)
  [[ "$tag" =~ ^v[1-9][0-9]*\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$ ]]
  gh api "repos/microsoft/TypeScript/releases/tags/$tag" \
    --jq 'select(.prerelease == false and .draft == false) | .tag_name' | grep -Fx "$tag"
  git fetch --no-tags stable-upstream "+refs/tags/$tag:refs/remotes/stable-upstream/release"
  bash "$trusted_merge" origin/main stable-upstream/release
  [[ "$(go run ./scripts/export-go.go --version)" == "$tag" ]]
  go run ./scripts/export-go.go
  git add --all -- public
  if ! git diff --cached --quiet; then
    git commit -m "Regenerate public module for TypeScript $tag"
  fi
  git push origin "$branch"
  if [[ $(gh pr view "$number" --repo "$GITHUB_REPOSITORY" --json isDraft --jq .isDraft) == true ]]; then
    gh pr ready "$number" --repo "$GITHUB_REPOSITORY"
  fi
  # The built-in token does not trigger pull_request workflows on bot pushes.
  gh workflow run compiler-packages.yml --repo "$GITHUB_REPOSITORY" --ref "$branch"
done < /tmp/typescript-update-prs
