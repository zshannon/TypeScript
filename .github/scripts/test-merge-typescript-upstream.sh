#!/usr/bin/env bash

set -euo pipefail

script_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
sync_script="$script_dir/merge-typescript-upstream.sh"
test_root=$(mktemp -d "${TMPDIR:-/tmp}/typescript-upstream-sync.XXXXXX")
trap 'rm -rf "$test_root"' EXIT

configure_git() {
  git config user.name "Upstream Sync Test"
  git config user.email "upstream-sync-test@example.invalid"
}

new_repo() {
  local path=$1
  git init -q -b main "$path"
  (
    cd "$path"
    configure_git
  )
}

assert_eq() {
  local expected=$1
  local actual=$2
  local message=$3
  if [[ $expected != "$actual" ]]; then
    echo "FAIL: $message (expected '$expected', got '$actual')" >&2
    exit 1
  fi
}

assert_contains() {
  local needle=$1
  local file=$2
  local message=$3
  if ! grep -Fq "$needle" "$file"; then
    echo "FAIL: $message" >&2
    echo "Expected to find '$needle' in $file" >&2
    exit 1
  fi
}

assert_not_exists() {
  local path=$1
  local message=$2
  if [[ -e $path ]]; then
    echo "FAIL: $message ($path exists)" >&2
    exit 1
  fi
}

test_existing_automation_branch_merge() {
  local repo="$test_root/clean"
  new_repo "$repo"
  (
    cd "$repo"
    printf 'shared base\n' >compiler.txt
    git add compiler.txt
    git commit -q -m "upstream base"
    base=$(git rev-parse HEAD)

    git switch -q -c fork-main
    printf 'fork integration\n' >fork.txt
    git add fork.txt
    git commit -q -m "fork change"

    git switch -q -c automation
    printf 'human review fix\n' >human.txt
    git add human.txt
    git commit -q -m "human automation-branch fix"
    human_commit=$(git rev-parse HEAD)

    git switch -q fork-main
    printf 'new target change\n' >target.txt
    git add target.txt
    git commit -q -m "target advanced"
    target_commit=$(git rev-parse HEAD)

    git switch -q main
    printf 'upstream feature\n' >upstream.txt
    git add upstream.txt
    git commit -q -m "upstream change"
    upstream_commit=$(git rev-parse HEAD)

    git switch -q automation
    output="$repo/output.txt"
    GITHUB_OUTPUT="$output" "$sync_script" fork-main main >/dev/null

    git merge-base --is-ancestor "$base" HEAD
    git merge-base --is-ancestor "$human_commit" HEAD
    git merge-base --is-ancestor "$target_commit" HEAD
    git merge-base --is-ancestor "$upstream_commit" HEAD
    assert_eq "3" "$(git rev-list --parents -n 1 HEAD | awk '{print NF}')" \
      "upstream synchronization must create a two-parent merge commit"
    assert_contains "fork integration" fork.txt "fork changes must survive the merge"
    assert_contains "human review fix" human.txt "human automation-branch fixes must survive the merge"
    assert_contains "new target change" target.txt "new target changes must be merged forward"
    assert_contains "upstream feature" upstream.txt "upstream changes must be present"
    assert_contains "upstream_changed=true" "$output" "clean merge must report an upstream change"
    assert_contains "branch_changed=true" "$output" "clean merge must report a branch change"
  )
  echo "PASS: existing automation branch preserves human fixes and merges target plus upstream"
}

test_no_update() {
  local repo="$test_root/no-update"
  new_repo "$repo"
  (
    cd "$repo"
    printf 'current upstream\n' >compiler.txt
    git add compiler.txt
    git commit -q -m "current upstream"
    git branch target
    git switch -q -c automation target
    before=$(git rev-parse HEAD)
    output="$repo/output.txt"

    GITHUB_OUTPUT="$output" "$sync_script" target main >/dev/null

    assert_eq "$before" "$(git rev-parse HEAD)" "an up-to-date run must not create a commit"
    assert_contains "upstream_changed=false" "$output" "up-to-date run must report no upstream change"
    assert_contains "branch_changed=false" "$output" "up-to-date run must report no branch change"
    assert_contains "ahead_count=0" "$output" "up-to-date run must have no PR commits"
  )
  echo "PASS: up-to-date synchronization is a no-op"
}

test_conflict_safety() {
  local repo="$test_root/conflict"
  new_repo "$repo"
  (
    cd "$repo"
    printf 'base\n' >compiler.txt
    git add compiler.txt
    git commit -q -m "common base"

    git switch -q -c fork-main
    printf 'fork edit\n' >compiler.txt
    git add compiler.txt
    git commit -q -m "fork edit"
    fork_head=$(git rev-parse HEAD)

    git switch -q main
    printf 'upstream edit\n' >compiler.txt
    git add compiler.txt
    git commit -q -m "upstream edit"

    git switch -q -c automation fork-main
    set +e
    stdout_file="$test_root/conflict.stdout.txt"
    stderr_file="$test_root/conflict.stderr.txt"
    "$sync_script" fork-main main >"$stdout_file" 2>"$stderr_file"
    status=$?
    set -e

    assert_eq "2" "$status" "a merge conflict must return the documented status"
    assert_eq "$fork_head" "$(git rev-parse HEAD)" "a conflict must leave the branch at its prior commit"
    assert_eq "fork edit" "$(cat compiler.txt)" "a conflict must preserve the fork's file"
    assert_eq "" "$(git status --porcelain)" "a conflict must be aborted to a clean worktree"
    assert_contains "conflicting files:" "$stderr_file" "conflict output must identify the failure"
    assert_contains "compiler.txt" "$stderr_file" "conflict output must name affected files"
  )
  echo "PASS: conflicts preserve fork commits and leave a clean worktree"
}

test_upstream_workflow_changes_are_excluded() {
  local repo="$test_root/workflow-tree"
  new_repo "$repo"
  (
    cd "$repo"
    mkdir -p .github/workflows
    printf 'fork edit baseline\n' >.github/workflows/edited.yml
    printf 'fork retained workflow\n' >.github/workflows/deleted-upstream.yml
    printf 'compiler base\n' >compiler.txt
    git add .github/workflows compiler.txt
    git commit -q -m "common base"
    git branch fork-main
    expected_tree=$(git rev-parse 'fork-main:.github/workflows')

    printf 'upstream replacement\n' >.github/workflows/edited.yml
    git rm -q .github/workflows/deleted-upstream.yml
    printf 'upstream-only workflow\n' >.github/workflows/added-upstream.yml
    git add .github/workflows
    git commit -q -m "upstream adds edits and deletes workflows"
    upstream_commit=$(git rev-parse HEAD)

    git switch -q -c automation fork-main
    "$sync_script" fork-main main >/dev/null

    assert_eq "$expected_tree" "$(git rev-parse 'HEAD:.github/workflows')" \
      "the complete workflow tree must match its pre-upstream state"
    assert_contains "fork edit baseline" .github/workflows/edited.yml \
      "an upstream workflow edit must be excluded"
    assert_contains "fork retained workflow" .github/workflows/deleted-upstream.yml \
      "an upstream workflow deletion must be excluded"
    assert_not_exists .github/workflows/added-upstream.yml \
      "an upstream workflow addition must be excluded"
    git merge-base --is-ancestor "$upstream_commit" HEAD
    assert_eq "3" "$(git rev-list --parents -n 1 HEAD | awk '{print NF}')" \
      "workflow-only upstream changes must still produce an ancestry merge"
    assert_eq "" "$(git status --porcelain)" "workflow restoration must leave a clean worktree"

    merged_head=$(git rev-parse HEAD)
    output="$test_root/workflow-only-noop-output.txt"
    GITHUB_OUTPUT="$output" "$sync_script" fork-main main >/dev/null
    assert_eq "$merged_head" "$(git rev-parse HEAD)" \
      "a repeated workflow-only synchronization must not create another merge"
    assert_contains "upstream_changed=false" "$output" \
      "a repeated workflow-only synchronization must report upstream as contained"
    assert_contains "branch_changed=false" "$output" \
      "a repeated workflow-only synchronization must report no branch change"
  )
  echo "PASS: upstream workflow-only changes merge ancestry once and then no-op"
}

test_workflow_only_conflict_is_resolved() {
  local repo="$test_root/workflow-conflict"
  new_repo "$repo"
  (
    cd "$repo"
    mkdir -p .github/workflows
    printf 'shared workflow\n' >.github/workflows/ci.yml
    git add .github/workflows/ci.yml
    git commit -q -m "common base"
    git branch fork-main

    git switch -q -c automation fork-main
    printf 'fork workflow\n' >.github/workflows/ci.yml
    git add .github/workflows/ci.yml
    git commit -q -m "fork workflow maintenance"
    fork_head=$(git rev-parse HEAD)

    git switch -q main
    printf 'upstream workflow\n' >.github/workflows/ci.yml
    git add .github/workflows/ci.yml
    git commit -q -m "upstream workflow maintenance"
    upstream_commit=$(git rev-parse HEAD)

    git switch -q automation
    "$sync_script" fork-main main >/dev/null

    assert_contains "fork workflow" .github/workflows/ci.yml \
      "a workflow-only conflict must retain the fork version"
    git merge-base --is-ancestor "$fork_head" HEAD
    git merge-base --is-ancestor "$upstream_commit" HEAD
    assert_eq "3" "$(git rev-list --parents -n 1 HEAD | awk '{print NF}')" \
      "a resolved workflow-only conflict must retain upstream ancestry"
    assert_eq "" "$(git status --porcelain)" "resolved workflow conflict must leave a clean worktree"
  )
  echo "PASS: workflow-only conflicts retain fork workflows and complete the merge"
}

test_mixed_workflow_and_source_conflicts_abort() {
  local repo="$test_root/mixed-conflict"
  new_repo "$repo"
  (
    cd "$repo"
    mkdir -p .github/workflows
    printf 'shared workflow\n' >.github/workflows/ci.yml
    printf 'shared compiler\n' >compiler.txt
    git add .github/workflows/ci.yml compiler.txt
    git commit -q -m "common base"
    git branch fork-main

    git switch -q -c automation fork-main
    printf 'fork workflow\n' >.github/workflows/ci.yml
    printf 'fork compiler\n' >compiler.txt
    git add .github/workflows/ci.yml compiler.txt
    git commit -q -m "fork changes"
    fork_head=$(git rev-parse HEAD)

    git switch -q main
    printf 'upstream workflow\n' >.github/workflows/ci.yml
    printf 'upstream compiler\n' >compiler.txt
    git add .github/workflows/ci.yml compiler.txt
    git commit -q -m "upstream changes"

    git switch -q automation
    stdout_file="$test_root/mixed-conflict.stdout.txt"
    stderr_file="$test_root/mixed-conflict.stderr.txt"
    set +e
    "$sync_script" fork-main main >"$stdout_file" 2>"$stderr_file"
    status=$?
    set -e

    assert_eq "2" "$status" "a mixed workflow and source conflict must fail"
    assert_eq "$fork_head" "$(git rev-parse HEAD)" "a mixed conflict must not create a merge commit"
    assert_eq "fork workflow" "$(cat .github/workflows/ci.yml)" \
      "a mixed conflict abort must restore the fork workflow"
    assert_eq "fork compiler" "$(cat compiler.txt)" \
      "a mixed conflict abort must restore the fork source"
    assert_contains ".github/workflows/ci.yml" "$stderr_file" \
      "mixed conflict output must identify the workflow conflict"
    assert_contains "compiler.txt" "$stderr_file" \
      "mixed conflict output must identify the source conflict"
    assert_eq "" "$(git status --porcelain)" "a mixed conflict abort must leave a clean worktree"
  )
  echo "PASS: mixed workflow and source conflicts abort without changing fork work"
}

test_target_workflow_maintenance_is_retained() {
  local repo="$test_root/target-workflows"
  new_repo "$repo"
  (
    cd "$repo"
    mkdir -p .github/workflows
    printf 'shared workflow\n' >.github/workflows/ci.yml
    printf 'compiler base\n' >compiler.txt
    git add .github/workflows/ci.yml compiler.txt
    git commit -q -m "common base"
    git branch fork-main

    git switch -q -c automation fork-main
    printf 'human automation workflow\n' >.github/workflows/human.yml
    git add .github/workflows/human.yml
    git commit -q -m "human workflow fix"
    human_commit=$(git rev-parse HEAD)

    git switch -q fork-main
    printf 'target-maintained workflow\n' >.github/workflows/ci.yml
    printf 'target-only workflow\n' >.github/workflows/target.yml
    git add .github/workflows
    git commit -q -m "target workflow maintenance"
    target_commit=$(git rev-parse HEAD)

    git switch -q main
    printf 'upstream workflow\n' >.github/workflows/ci.yml
    printf 'upstream-only workflow\n' >.github/workflows/upstream.yml
    printf 'upstream compiler\n' >upstream.txt
    git add .github/workflows upstream.txt
    git commit -q -m "upstream workflow and source changes"
    upstream_commit=$(git rev-parse HEAD)

    git switch -q automation
    "$sync_script" fork-main main >/dev/null

    assert_contains "target-maintained workflow" .github/workflows/ci.yml \
      "target-side workflow edits must survive upstream merge"
    assert_contains "target-only workflow" .github/workflows/target.yml \
      "target-side workflow additions must survive upstream merge"
    assert_contains "human automation workflow" .github/workflows/human.yml \
      "human automation-branch workflows must survive upstream merge"
    assert_not_exists .github/workflows/upstream.yml \
      "upstream-only workflows must not enter the fork"
    assert_contains "upstream compiler" upstream.txt "non-workflow upstream changes must be merged"
    git merge-base --is-ancestor "$human_commit" HEAD
    git merge-base --is-ancestor "$target_commit" HEAD
    git merge-base --is-ancestor "$upstream_commit" HEAD
    assert_eq "" "$(git status --porcelain)" "target workflow preservation must leave a clean worktree"
  )
  echo "PASS: target and human workflow maintenance survive the upstream merge"
}

test_absent_workflow_tree_remains_absent() {
  local repo="$test_root/no-workflow-tree"
  new_repo "$repo"
  (
    cd "$repo"
    printf 'compiler base\n' >compiler.txt
    git add compiler.txt
    git commit -q -m "common base"
    git branch fork-main

    mkdir -p .github/workflows
    printf 'upstream-only workflow\n' >.github/workflows/ci.yml
    git add .github/workflows/ci.yml
    git commit -q -m "upstream adds first workflow"
    upstream_commit=$(git rev-parse HEAD)

    git switch -q -c automation fork-main
    mkdir -p .github/workflows
    "$sync_script" fork-main main >/dev/null

    assert_not_exists .github/workflows/ci.yml \
      "an upstream workflow must not populate a previously empty workflow tree"
    git merge-base --is-ancestor "$upstream_commit" HEAD
    assert_eq "" "$(git status --porcelain)" "empty workflow preservation must leave a clean worktree"
  )
  echo "PASS: a missing workflow tree stays empty while upstream ancestry is merged"
}

test_generated_public_conflict_is_discarded() {
  local repo="$test_root/public-conflict"
  new_repo "$repo"
  (
    cd "$repo"
    mkdir -p public/compiler
    printf 'base generated compiler\n' >public/compiler/generated.go
    printf 'compiler base\n' >compiler.txt
    git add compiler.txt public/compiler/generated.go
    git commit -q -m "common base"
    git branch fork-main

    git switch -q -c automation fork-main
    printf 'automation generated compiler\n' >public/compiler/generated.go
    git add public/compiler/generated.go
    git commit -q -m "automation generated output"
    automation_head=$(git rev-parse HEAD)

    git switch -q fork-main
    printf 'target generated compiler\n' >public/compiler/generated.go
    printf 'target compiler change\n' >target.txt
    git add public/compiler/generated.go target.txt
    git commit -q -m "target source and generated output"
    target_head=$(git rev-parse HEAD)

    git switch -q main
    printf 'upstream compiler change\n' >upstream.txt
    git add upstream.txt
    git commit -q -m "upstream source change"
    upstream_head=$(git rev-parse HEAD)

    git switch -q automation
    "$sync_script" fork-main main >/dev/null

    assert_not_exists public \
      "conflicted generated output must be removed for post-merge regeneration"
    assert_contains "target compiler change" target.txt \
      "target source changes must survive a generated-only conflict"
    assert_contains "upstream compiler change" upstream.txt \
      "upstream source changes must merge after a generated-only conflict"
    git merge-base --is-ancestor "$automation_head" HEAD
    git merge-base --is-ancestor "$target_head" HEAD
    git merge-base --is-ancestor "$upstream_head" HEAD
    assert_eq "" "$(git status --porcelain)" \
      "generated-only conflict resolution must leave a clean worktree"
  )
  echo "PASS: generated public conflicts are discarded for post-merge regeneration"
}

test_existing_automation_branch_merge
test_no_update
test_conflict_safety
test_upstream_workflow_changes_are_excluded
test_workflow_only_conflict_is_resolved
test_mixed_workflow_and_source_conflicts_abort
test_target_workflow_maintenance_is_retained
test_absent_workflow_tree_remains_absent
test_generated_public_conflict_is_discarded
echo "All TypeScript upstream merge integration tests passed."
