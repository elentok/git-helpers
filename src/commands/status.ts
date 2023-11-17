import { updateRemote } from "../lib/git.ts"
import { CHECKMARK, ERROR } from "../lib/helpers.ts"
import { findRepoOrExit } from "../lib/repo.ts"
import { getStatus } from "../lib/status.ts"

export function status({ quick }: { quick?: boolean } = {}) {
  const repo = findRepoOrExit(Deno.cwd())
  if (!quick) {
    updateRemote(repo)
  }
  const status = getStatus(repo)

  for (const localBranch of status.localBranches) {
    if (localBranch.remoteBranches.length === 0) {
      console.info(`- ${localBranch.name} (local only)`)
    } else {
      const symbol = localBranch.isSynced ? CHECKMARK : ERROR

      if (localBranch.remoteBranches.length === 1) {
        const rb = localBranch.remoteBranches[0]
        console.info(`${symbol} ${localBranch.gitName} (${rb.status.pretty})`)
      } else {
        console.info(`${symbol} ${localBranch.gitName}:`)
        for (const rb of localBranch.remoteBranches) {
          const symbol = rb.status.name === "same" ? CHECKMARK : ERROR
          console.info(`  ${symbol} ${rb.gitName} (${rb.status.pretty})`)
        }
      }
    }
  }
}
