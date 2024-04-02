import { isPresent } from "./helpers.ts"
import { Repo } from "./Repo.ts"
import { Worktree } from "./types.ts"

export class RepoWorktree {
  constructor(private repo: Repo) {}

  list(): Worktree[] {
    if (!this.repo.isBare) {
      throw new Error("Only bare repositories can have worktrees")
    }

    const { stdout } = this.repo.git(["worktree", "list"])
    return stdout.split("\n").map((line) => {
      const match = /^([^\s]+)\s+[^\s]+\s+\[(.*)\]/.exec(line)
      if (match == null) {
        return
      }

      const fullPath = match[1]
      const name = fullPath.substring(this.repo.root.length + 1)

      return { fullPath, name, branchName: match[2] }
    }).filter(isPresent)
  }

  remove(name: string, { force = false }: { force?: boolean } = {}): void {
    this.repo.git(
      ["worktree", "remove", force ? "-f" : undefined, name].filter(isPresent),
    )
  }
}
