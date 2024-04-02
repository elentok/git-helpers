import { LocalBranch, RemoteBranch, Repo } from "./git/types.ts"
import * as git from "./git/index.ts"

export interface RepoStatus {
  localBranches: LocalBranchStatus[]
  hasUncommitedChanges?: boolean
  hasUntrackedFiles?: boolean
}

export interface LocalBranchStatus extends LocalBranch {
  remoteBranches: RemoteBranchStatus[]
  isSynced: boolean
  hash: string
}

export interface RemoteBranchStatus extends RemoteBranch {
  status: SyncStatus
  hash: string
}

export interface SyncStatus {
  behind: number
  ahead: number
  name: "behind" | "ahead" | "diverged" | "same" | "unclear"
  pretty: string
}

export function getRepoStatus(repo: Repo): RepoStatus {
  const branches = git.branch.list(repo)

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

  const isBare = git.isBare(repo)

  return {
    localBranches: localBranchStatuses,
    hasUncommitedChanges: isBare ? undefined : git.hasUncommitedChanges(repo),
    hasUntrackedFiles: isBare ? undefined : git.hasUntrackedFiles(repo),
  }
}

function getBranchStatus(
  repo: Repo,
  localBranch: LocalBranch,
  remoteBranches: RemoteBranch[],
): LocalBranchStatus {
  const localHash = git.hash(repo, localBranch.gitName)

  const remoteBranchStatuses: RemoteBranchStatus[] = remoteBranches.map(
    (rb) => {
      const remoteHash = git.hash(repo, rb.gitName)
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

  const behind = git.revCount(repo, localBranch.gitName, remoteBranch.gitName)
  const ahead = git.revCount(repo, remoteBranch.gitName, localBranch.gitName)

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
