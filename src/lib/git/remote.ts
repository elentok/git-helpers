import { run } from "./run.ts"
import { Repo } from "./types.ts"

export function list(repo: string | Repo): string[] {
  return run(repo, ["remote"]).stdout.split("\n")
}

export function update(repo: string | Repo): void {
  console.info("Updating remotes...")
  run(repo, ["remote", "update"])
}

export function pruneAll(repo: string | Repo): void {
  for (const remote of list(repo)) {
    prune(repo, remote)
  }
}

export function prune(repo: string | Repo, remote: string): void {
  console.info(`Pruning remote ${remote}...`)
  run(repo, ["remote", "prune", remote])
}
