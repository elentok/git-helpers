#!/usr/bin/env bash

set -euo pipefail

TEST_DIR=~/tmp/git-helpers-test

function main() {
  rm -rf $TEST_DIR
  mkdir -p $TEST_DIR
  cd $TEST_DIR

  create-normal-repo
  create-worktree-repo
}

function create-normal-repo() {
  cd "$TEST_DIR"

  mkdir bare
  cd bare
  git init --bare

  cd ..
  git clone bare clone1
  git clone bare clone2

  cd clone1

  add-and-commit file1.txt
  git push origin main

  add-and-commit file2.txt
  git checkout -b branch2
  git push origin branch2
  git checkout main

  cd ../clone2
  git pull origin main

  add-and-commit file3.txt
  git checkout -b branch3
  git push origin branch3
  git checkout main
}

function create-worktree-repo() {
  cd "$TEST_DIR"

  mkdir worktrees
  cd worktrees
  git init --bare
  git remote add origin "$TEST_DIR/bare"
  echo "Remote update"
  git remote update

  git worktree add --checkout branch2 origin/branch2
  git worktree add --checkout branch3 origin/branch3
}

function add-and-commit() {
  echo "this is $1" > "$1"
  git add "$1"
  git commit -m "commit $1"
}

main "$@"
