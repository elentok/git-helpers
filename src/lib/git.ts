import { ShellOptions, ShellResult, shell } from "./shell.ts"
import { Repo } from "./types.ts"

export function git(repo: Repo | string, args: string[], options?: ShellOptions): ShellResult {
  const cwd = typeof repo === "string" ? repo : repo.root
  return shell("git", { args, cwd, ...options })
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

export function hasUncommitedChanges(repo: Repo): boolean {
  const output = git(repo, ["status", "--porcelain=v1"]).stdout
  if (output.length === 0) return false

  const lines = output.split("\n").filter((l) => l.length > 0 && !l.startsWith("?? "))
  return lines.length > 0
}

export function hasUntrackedFiles(repo: Repo): boolean {
  const output = git(repo, ["status", "--porcelain=v1"]).stdout
  if (output.length === 0) return false

  const lines = output.split("\n")
  return lines.find((l) => l.startsWith("?? ")) != null
}

export function revParseString(repo: Repo | string, what: "show-toplevel"): string | undefined {
  const result = git(repo, ["rev-parse", `--${what}`], { throwError: false })
  if (!result.success) return
  return result.stdout
}

export function revParseBoolean(
  repo: Repo | string,
  what: "is-bare-repository" | "is-inside-work-tree",
): boolean {
  const result = git(repo, ["rev-parse", `--${what}`], { throwError: false })
  return result.success && result.stdout === "true"
}
