import { ShellOptions, ShellResult, shell } from "./shell.ts"
import { Repo } from "./types.ts"

export function git(repo: Repo, args: string[], options?: ShellOptions): ShellResult {
  return shell("git", { args, cwd: repo.root, ...options })
}

export function gitRemotes(repo: Repo): string[] {
  return git(repo, ["remote"]).stdout.split("\n")
}
