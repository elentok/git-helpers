import { findRepoOrThrow } from "./findRepo.ts"
import { run } from "./run.ts"
import { Repo } from "./types.ts"

export function hash(repo: string | Repo, ref: string): string {
  return run(repo, ["log", "-1", "--pretty=%H", ref]).stdout
}

export function revCount(
  repo: string | Repo,
  fromRef: string,
  toRef: string,
): number {
  const output =
    run(repo, ["rev-list", "--count", `${fromRef}..${toRef}`]).stdout
  const count = Number(output)

  if (isNaN(count)) {
    throw new Error(`Invalid rev-list count '${output}'`)
  }

  return count
}

export function hasUncommitedChanges(repo: string | Repo): boolean {
  const output = run(repo, ["status", "--porcelain=v1"]).stdout
  if (output.length === 0) return false

  const lines = output.split("\n").filter((l) =>
    l.length > 0 && !l.startsWith("?? ")
  )
  return lines.length > 0
}

export function hasUntrackedFiles(repo: string | Repo): boolean {
  const output = run(repo, ["status", "--porcelain=v1"]).stdout
  if (output.length === 0) return false

  const lines = output.split("\n")
  return lines.find((l) => l.startsWith("?? ")) != null
}

export function isBare(repo: string | Repo): boolean {
  const repoWithDetails = (typeof repo === "string")
    ? findRepoOrThrow(repo)
    : repo

  return repoWithDetails.isBare
}

export function root(repo: string | Repo): string {
  return (typeof repo === "string") ? repo : repo.root
}
