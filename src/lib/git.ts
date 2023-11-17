import { ShellOptions, ShellResult, shell } from "./shell.ts"
import { LocalBranch, LocalBranchStatus, Repo } from "./types.ts"

export function git(repo: Repo, args: string[], options?: ShellOptions): ShellResult {
  return shell("git", { args, cwd: repo.root, ...options })
}

export function getRemotes(repo: Repo): string[] {
  return git(repo, ["remote"]).stdout.split("\n")
}

export function updateRemote(repo: Repo): void {
  console.info("Updating remotes...")
  git(repo, ["remote", "update"])
}

export function getHash(repo: Repo, ref: string): string {
  return git(repo, ["log", "-1", "--pretty=%H", ref]).stdout
}

export function getCurrentBranch(repo: Repo): string {
  return git(repo, ["rev-parse", "--abbrev-ref", "HEAD"]).stdout
}

export function getRevCount(repo: Repo, fromRef: string, toRef: string): number {
  const output = git(repo, ["rev-list", "--count", `${fromRef}..${toRef}`]).stdout
  const count = Number(output)

  if (isNaN(count)) {
    throw new Error(`Invalid rev-list count '${output}'`)
  }

  return count
}
