import { run } from "./run.ts"
import { Repo } from "./types.ts"

export function list(repo: Repo): string[] {
  return run(repo, ["remote"]).stdout.split("\n")
}

export function update(repo: Repo): void {
  console.info("Updating remotes...")
  run(repo, ["remote", "update"])
}
