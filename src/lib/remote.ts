import { git } from "./git.ts"
import { Repo } from "./types.ts"

export function remotes(repo: Repo): string[] {
  return git(repo, ["remote"]).split("\n")
}
