import { isPresent } from "../helpers.ts"
import { run } from "./run.ts"
import { Repo, Worktree } from "./types.ts"
import { isBare, root } from "./utils.ts"

export function list(repo: string | Repo): Worktree[] {
  if (!isBare(repo)) {
    return []
  }

  const { stdout } = run(repo, ["worktree", "list"])
  return stdout.split("\n").map((line) => {
    const match = /^([^\s]+)\s+[^\s]+\s+\[(.*)\]/.exec(line)
    if (match == null) {
      return
    }

    const fullPath = match[1]
    const name = fullPath.substring(root(repo).length + 1)

    return { fullPath, name, branchName: match[2] }
  }).filter(isPresent)
}

export function remove(
  repo: string | Repo,
  name: string,
  { force = false }: { force?: boolean } = {},
): void {
  run(
    repo,
    ["worktree", "remove", force ? "-f" : undefined, name].filter(isPresent),
  )
}
