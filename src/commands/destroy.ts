import { findRepoOrExit } from "../lib/repo.ts"
import { pickBranch, pickRemote } from "../lib/pickers.ts"
import { gitBranches } from "../lib/branch.ts"
import { RemoteBranch, LocalBranch } from "../lib/types.ts"

export async function destroy() {
  console.log("[elentok] [destroy.ts] destroy")

  const repo = findRepoOrExit(Deno.cwd())
  console.log("[elentok] [destroy.ts] destroy", repo)

  // const branches = gitBranches(repo)
  // console.log("[elentok] [destroy.ts] destroy", branches)

  // const branches = await pickBranch(repo, { allowMultiple: true })
  // console.log("[elentok] [destroy.ts] destroy", branches)
  //
  // const localBanches: LocalBranch[] = []
  // const remoteBranches: RemoteBranch[] = []
  //
  // for (const branch of branches) {
  //   // if (branch.type === "remote")
  // }

  // const remote = await pickRemote(repo, { selectOne: true })
  // console.log("[elentok] [destroy.ts] destroy", remote)
}
