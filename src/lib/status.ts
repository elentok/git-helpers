import { gitBranches } from "./branch.ts"
import { getHash, getRevCount, git } from "./git.ts"
import {
  LocalBranch,
  LocalBranchStatus,
  RemoteBranch,
  RemoteBranchStatus,
  Repo,
  RepoStatus,
  SyncStatus,
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
    const status = getSyncStatus(repo, localBranch, localHash, rb, remoteHash)

    return { ...rb, hash: remoteHash, status }
  })

  const isSynced = !remoteBranchStatuses.find((rb) => rb.status.name !== "same")

  return { ...localBranch, isSynced, hash: localHash, remoteBranches: remoteBranchStatuses }
}

function getSyncStatus(
  repo: Repo,
  localBranch: LocalBranch,
  localHash: string,
  remoteBranch: RemoteBranch,
  remoteHash: string,
): SyncStatus {
  if (localHash === remoteHash) {
    return { name: "same", ahead: 0, behind: 0 }
  }

  const behind = getRevCount(repo, localBranch.gitName, remoteBranch.gitName)
  const ahead = getRevCount(repo, remoteBranch.gitName, localBranch.gitName)

  let name: SyncStatus["name"] = "unclear"
  if (ahead > 0) {
    if (behind > 0) {
      name = "diverged"
    } else {
      name = "ahead"
    }
  } else {
    name = "behind"
  }

  return { name, behind, ahead }
}
