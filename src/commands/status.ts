import { gitBranches } from "../lib/branch.ts"
import { getHash } from "../lib/git.ts"
import { findRepoOrExit } from "../lib/repo.ts"
import { Branch, LocalBranch, RemoteBranch, Repo } from "../lib/types.ts"

export function status(branch?: string) {
  console.log("[elentok] [status.ts] status", { branch })

  const repo = findRepoOrExit(Deno.cwd())
  const status = getStatus(repo)
  console.log("[elentok] [status.ts] status", status)
  // const branches = gitBranches(repo)

  // const localBranches = branches.filter((b) => b.type === "local")
  // const remoteBranches = branches.filter((b) => b.type === "remote")
  //
  // for (const branch of localBranches) {
  //   const matchingRemoteBranches =
  // }
}

function getStatus(repo: Repo): RepoStatus {
  const branches = gitBranches(repo)

  const localBranches = branches.filter((b) => b.type === "local") as LocalBranch[]
  const remoteBranches = branches.filter((b) => b.type === "remote") as RemoteBranch[]

  const detailedRemoteBranches = remoteBranches.map((branch) => ({
    ...branch,
    hash: getHash(repo, branch.gitName),
  }))

  const detailedLocalBranches = localBranches.map((branch) => ({
    ...branch,
    hash: getHash(repo, branch.gitName),
    remoteBranches: detailedRemoteBranches.filter((b) => b.name === branch.name),
  }))

  return {
    localBranches: detailedLocalBranches,
    remoteBranches: detailedRemoteBranches,
  }
}

interface RepoStatus {
  localBranches: DetailedLocalBranch[]
  remoteBranches: DetailedRemoteBranch[]
}

interface DetailedLocalBranch extends LocalBranch {
  hash: string
  remoteBranches: RemoteBranch[]
}

interface DetailedRemoteBranch extends RemoteBranch {
  hash: string
}
