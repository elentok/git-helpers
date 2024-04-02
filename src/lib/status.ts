import { gitBranches } from "./branch.ts"
import { Repo } from "./repo.ts"
import {
  LocalBranch,
  LocalBranchStatus,
  RemoteBranch,
  RemoteBranchStatus,
  RepoStatus,
  SyncStatus,
} from "./types.ts"

export function getRepoStatus(repo: Repo): RepoStatus {
  const branches = repo.branches()

  const localBranches = branches.filter((b) =>
    b.type === "local"
  ) as LocalBranch[]
  const remoteBranches = branches.filter((b) =>
    b.type === "remote"
  ) as RemoteBranch[]

  const localBranchStatuses = localBranches.map((localBranch) => {
    const relatedRemoteBranches = remoteBranches.filter((b) =>
      b.name === localBranch.name
    )
    return getBranchStatus(repo, localBranch, relatedRemoteBranches)
  })

  return {
    localBranches: localBranchStatuses,
    hasUncommitedChanges: repo.isBare ? undefined : repo.hasUncommitedChanges(),
    hasUntrackedFiles: repo.isBare ? undefined : repo.hasUntrackedFiles(),
  }
}

function getBranchStatus(
  repo: Repo,
  localBranch: LocalBranch,
  remoteBranches: RemoteBranch[],
): LocalBranchStatus {
  const localHash = repo.hash(localBranch.gitName)

  const remoteBranchStatuses: RemoteBranchStatus[] = remoteBranches.map(
    (rb) => {
      const remoteHash = repo.hash(rb.gitName)
      const status = getSyncStatus(repo, localBranch, localHash, rb, remoteHash)

      return { ...rb, hash: remoteHash, status }
    },
  )

  const isSynced = !remoteBranchStatuses.find((rb) => rb.status.name !== "same")

  return {
    ...localBranch,
    isSynced,
    hash: localHash,
    remoteBranches: remoteBranchStatuses,
  }
}

function getSyncStatus(
  repo: Repo,
  localBranch: LocalBranch,
  localHash: string,
  remoteBranch: RemoteBranch,
  remoteHash: string,
): SyncStatus {
  if (localHash === remoteHash) {
    return { name: "same", ahead: 0, behind: 0, pretty: "synced" }
  }

  const behind = repo.revCount(localBranch.gitName, remoteBranch.gitName)
  const ahead = repo.revCount(remoteBranch.gitName, localBranch.gitName)

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

  const pretty = [
    name === "ahead" || name === "behind" ? null : name,
    ahead > 0 ? `${ahead} ahead` : null,
    behind > 0 ? `${behind} behind` : null,
  ]
    .filter((d) => d != null)
    .join(", ")

  return { name, behind, ahead, pretty }
}
