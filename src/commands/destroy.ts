import { fzf } from "../lib/fzf.ts"
import * as git from "../lib/git/index.ts"
import { getRepoStatus } from "../lib/status.ts"

export async function destroy({ branchOnly }: { branchOnly?: boolean }) {
  const repo = git.findRepoOrExit(Deno.cwd())

  const worktrees = git.worktree.list(repo)
  const status = getRepoStatus(repo)

  const items = status.localBranches.map((b) => {
    const synced = b.isSynced ? "synced" : "not synced"
    const status = b.remoteBranches.map((rb) =>
      `[${rb.remoteName}: ${rb.status.pretty}]`
    ).join(", ")
    return `${b.name} (${synced}) ${status}`
  })

  const selectedItems = await fzf({
    items,
    allowMultiple: true,
  })

  for (const item of selectedItems) {
    const branchName = item.split(" ")[0]
    const branch = status.localBranches.find((b) => b.name === branchName)
    if (branch == null) {
      console.error(`Invalid branch name: "${branchName}"`)
      continue
    }
    for (const remoteBranch of branch.remoteBranches) {
      if (remoteBranch.status.name === "same") {
        console.info(`- Deleting remote branch ${remoteBranch.gitName}`)
        git.branch.deleteRemoteBranch(repo, remoteBranch)
      }
    }

    if (!branchOnly) {
      const worktree = worktrees.find((w) => w.branchName === branchName)
      if (worktree != null) {
        console.info(`- Deleting worktree ${worktree.name}`)
        git.worktree.remove(repo, worktree.name, { force: true })
      }
    }

    console.info(`- Deleting local branch ${branch.gitName}`)
    git.branch.deleteLocalBranch(repo, branch.name, { force: true })
  }
}
