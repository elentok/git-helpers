import { findRepoOrExit } from "../lib/repo.ts"
import { pickBranch, pickRemote } from "../lib/pickers.ts"
import { gitBranches } from "../lib/branch.ts"

export async function destroy() {
  console.log("[elentok] [destroy.ts] destroy")

  const repo = findRepoOrExit(Deno.cwd())
  console.log("[elentok] [destroy.ts] destroy", repo)

  // const branches = gitBranches(repo)
  // console.log("[elentok] [destroy.ts] destroy", branches)

  const branches = await pickBranch(repo)
  console.log("[elentok] [destroy.ts] destroy", branches)

  // const remote = await pickRemote(repo, { selectOne: true })
  // console.log("[elentok] [destroy.ts] destroy", remote)
}
