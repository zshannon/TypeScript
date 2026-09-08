#!/usr/bin/env bash

set -euo pipefail

if [[ $# -ne 2 ]]; then
  echo "usage: $0 <target-ref> <upstream-ref>" >&2
  exit 64
fi

target_ref=$1
upstream_ref=$2

for ref in "$target_ref" "$upstream_ref"; do
  if ! git rev-parse --verify "${ref}^{commit}" >/dev/null 2>&1; then
    echo "error: $ref does not resolve to a commit" >&2
    exit 64
  fi
done

if [[ -n $(git status --porcelain --untracked-files=normal) ]]; then
  echo "error: refusing to sync from a dirty worktree" >&2
  git status --short >&2
  exit 65
fi

initial_head=$(git rev-parse HEAD)
upstream_sha=$(git rev-parse "${upstream_ref}^{commit}")
upstream_changed=false

write_output() {
  local key=$1
  local value=$2
  if [[ -n ${GITHUB_OUTPUT:-} ]]; then
    printf '%s=%s\n' "$key" "$value" >>"$GITHUB_OUTPUT"
  fi
  printf '%s=%s\n' "$key" "$value"
}

report_merge_failure() {
  local label=$1
  local ref=$2
  local conflicts
  conflicts=$(git diff --name-only --diff-filter=U)

  echo "error: could not merge $label ($ref); the automation branch was not pushed" >&2
  if [[ -n $conflicts ]]; then
    echo "conflicting files:" >&2
    printf '  %s\n' "$conflicts" >&2
  else
    echo "git rejected the merge before creating a merge result; inspect the log above" >&2
  fi

  if [[ -n ${GITHUB_STEP_SUMMARY:-} ]]; then
    {
      echo "## TypeScript upstream merge needs human resolution"
      echo
      echo "The sync runner could not merge \`$label\` (\`$ref\`). Nothing was pushed, so existing commits on the automation branch are unchanged."
      if [[ -n $conflicts ]]; then
        echo
        echo "Conflicting files:"
        echo '```text'
        printf '%s\n' "$conflicts"
        echo '```'
      fi
    } >>"$GITHUB_STEP_SUMMARY"
  fi

  if git rev-parse --verify -q MERGE_HEAD >/dev/null; then
    git merge --abort
  fi
}

merge_ref() {
  local label=$1
  local ref=$2
  local conflicts
  local non_public_conflicts=""
  shift 2

  if git merge-base --is-ancestor "$ref" HEAD; then
    echo "$label is already contained in the automation branch"
    return 0
  fi

  set +e
  git merge "$@" "$ref"
  local merge_status=$?
  set -e
  if ((merge_status != 0)); then
    conflicts=$(git diff --name-only --diff-filter=U)
    while IFS= read -r conflict; do
      [[ -z $conflict ]] && continue
      case "$conflict" in
        public/*) ;;
        *) non_public_conflicts+="${non_public_conflicts:+$'\n'}$conflict" ;;
      esac
    done <<<"$conflicts"

    if [[ -n $conflicts && -z $non_public_conflicts ]]; then
      echo "Resolving generated public module conflicts by regenerating after the merge"
      if git rm -r -q -f --ignore-unmatch -- public && git commit --no-edit; then
        return 0
      fi
    fi

    report_merge_failure "$label" "$ref"
    return 2
  fi
}

restore_workflow_tree() {
  local treeish=$1

  # Remove every workflow introduced or changed by upstream, including unmerged
  # entries, before restoring the exact tree owned by the fork.
  if ! git rm -r -q -f --ignore-unmatch -- .github/workflows; then
    echo "error: failed to remove the upstream workflow merge result" >&2
    return 1
  fi
  if git cat-file -e "${treeish}:.github/workflows" 2>/dev/null; then
    if ! git restore --source="$treeish" --staged --worktree -- .github/workflows; then
      echo "error: failed to restore the fork-owned workflow tree" >&2
      return 1
    fi
  fi
}

merge_upstream() {
  local ref=$1
  local pre_upstream_head
  local merge_status
  local conflicts
  local non_workflow_conflicts=""
  local has_public_conflicts=false

  if git merge-base --is-ancestor "$ref" HEAD; then
    echo "Microsoft TypeScript upstream is already contained in the automation branch"
    return 0
  fi

  upstream_changed=true
  pre_upstream_head=$(git rev-parse HEAD)

  set +e
  git merge --no-commit --no-ff -m "Merge Microsoft TypeScript stable release" "$ref"
  merge_status=$?
  set -e

  if ((merge_status != 0)); then
    conflicts=$(git diff --name-only --diff-filter=U)
    if [[ -z $conflicts ]]; then
      report_merge_failure "Microsoft TypeScript upstream" "$ref"
      return 2
    fi

    while IFS= read -r conflict; do
      [[ -z $conflict ]] && continue
      case "$conflict" in
        .github/workflows/*) ;;
        public/*) has_public_conflicts=true ;;
        *) non_workflow_conflicts+="${non_workflow_conflicts:+$'\n'}$conflict" ;;
      esac
    done <<<"$conflicts"

    if [[ -n $non_workflow_conflicts ]]; then
      report_merge_failure "Microsoft TypeScript upstream" "$ref"
      return 2
    fi

    echo "Resolving fork-owned workflows and generated public module conflicts"
  fi

  if ! restore_workflow_tree "$pre_upstream_head"; then
    report_merge_failure "Microsoft TypeScript upstream" "$ref"
    return 2
  fi

  if [[ $has_public_conflicts == true ]]; then
    if ! git rm -r -q -f --ignore-unmatch -- public; then
      echo "error: failed to remove conflicted generated public module" >&2
      report_merge_failure "Microsoft TypeScript upstream" "$ref"
      return 2
    fi
  fi

  conflicts=$(git diff --name-only --diff-filter=U)
  if [[ -n $conflicts ]]; then
    report_merge_failure "Microsoft TypeScript upstream" "$ref"
    return 2
  fi

  if ! git diff --quiet "$pre_upstream_head" -- .github/workflows; then
    echo "error: failed to restore the fork-owned workflow tree" >&2
    report_merge_failure "Microsoft TypeScript upstream" "$ref"
    return 2
  fi

  if ! git commit -m "Merge Microsoft TypeScript stable release"; then
    echo "error: failed to commit the prepared upstream merge" >&2
    report_merge_failure "Microsoft TypeScript upstream" "$ref"
    return 2
  fi
}

merge_ref "target branch" "$target_ref" --no-edit -m "Merge target branch before TypeScript upstream sync" || exit $?

if ! git merge-base --is-ancestor "$upstream_ref" HEAD; then
  merge_upstream "$upstream_ref" || exit $?
else
  echo "Microsoft TypeScript upstream is already contained in the automation branch"
fi

final_head=$(git rev-parse HEAD)
ahead_count=$(git rev-list --count "${target_ref}..HEAD")
branch_changed=false
if [[ $final_head != "$initial_head" ]]; then
  branch_changed=true
fi

write_output upstream_sha "$upstream_sha"
write_output upstream_changed "$upstream_changed"
write_output branch_changed "$branch_changed"
write_output ahead_count "$ahead_count"
