import { updateRemote } from "../lib/git.ts"
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
      const symbol = localBranch.isSynced ? "v" : "x"
      const suffix = localBranch.remoteBranches.length > 1 ? ":" : ""

      console.info(`${symbol} ${localBranch.name}${suffix}`)
      if (localBranch.remoteBranches.length > 1) {
        for (const remoteBranch of localBranch.remoteBranches) {
          const symbol = remoteBranch.status === "same" ? "v" : "x"
          console.info(`${symbol} ${localBranch.name}`)
        }
      }
    }
  }
}
