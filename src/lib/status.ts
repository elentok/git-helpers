import { gitBranches } from "./branch.ts"
import { getHash } from "./git.ts"
import {
  LocalBranch,
  LocalBranchStatus,
  RemoteBranch,
  RemoteBranchStatus,
  Repo,
  RepoStatus,
} from "./types.ts"

export function getStatus(repo: Repo): RepoStatus {
  const branches = gitBranches(repo)

  const localBranches = branches.filter((b) => b.type === "local") as LocalBranch[]
  const remoteBranches = branches.filter((b) => b.type === "remote") as RemoteBranch[]

  const localBranchStatuses = localBranches.map((localBranch) => {
    const relatedRemoteBranches = remoteBranches.filter((b) => b.name === localBranch.name)
    return getBranchStatus(repo, localBranch, relatedRemoteBranches)
  })

  return {
    localBranches: localBranchStatuses,
  }
}

function getBranchStatus(
  repo: Repo,
  localBranch: LocalBranch,
  remoteBranches: RemoteBranch[],
): LocalBranchStatus {
  const localHash = getHash(repo, localBranch.gitName)

  const remoteBranchStatuses: RemoteBranchStatus[] = remoteBranches.map((rb) => {
    const remoteHash = getHash(repo, rb.gitName)
    return { ...rb, hash: remoteHash, status: localHash === remoteHash ? "same" : "unclear" }
  })

  const isSynced = !remoteBranchStatuses.find((rb) => rb.status !== "same")

  return { ...localBranch, isSynced, hash: localHash, remoteBranches: remoteBranchStatuses }
}
