import { ShellOptions, ShellResult, shell } from "./shell.ts"
import { Repo } from "./types.ts"

export function git(repo: Repo, args: string[], options?: ShellOptions): ShellResult {
  return shell("git", { args, cwd: repo.root, ...options })
}

export function getRemotes(repo: Repo): string[] {
  return git(repo, ["remote"]).stdout.split("\n")
}

export function getHash(repo: Repo, ref: string): string {
  return git(repo, ["log", "-1", "--pretty=%H", ref]).stdout
}

export function getCurrentBranch(repo: Repo): string {
  return git(repo, ["rev-parse", "--abbrev-ref", "HEAD"]).stdout
}
