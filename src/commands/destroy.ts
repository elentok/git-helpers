import { findRepoOrExit } from "../lib/repo.ts"
import { fzf } from "../lib/fzf.ts"
import { getRepoStatus } from "../lib/status.ts"

export async function destroy() {
  const repo = findRepoOrExit(Deno.cwd())

  const status = getRepoStatus(repo)

  const items = status.localBranches.map((b) =>
    `${b.name} (${b.isSynced ? "synced" : "not synced"})`
  )

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
        repo.deleteRemoteBranch(remoteBranch)
      }
    }

    console.info(`- Deleting local branch ${branch.gitName}`)
    repo.deleteLocalBranch(branch.name, { force: !branch.isSynced })
  }
}
