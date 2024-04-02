import { run } from "./run.ts"
import { Repo } from "./types.ts"

export function list(repo: string | Repo): string[] {
  return run(repo, ["remote"]).stdout.split("\n")
}

export function update(repo: string | Repo): void {
  console.info("Updating remotes...")
  run(repo, ["remote", "update"])
}
