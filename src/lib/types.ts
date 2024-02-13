export interface Remote {
  name: string
  url: string
}

export class ShellError extends Error {
  constructor(
    public code: number,
    public command: string,
    public output: string,
  ) {
    super(
      `Shell command '${command}' failed with exitcode ${code}:\n\n${output}`,
    )
  }
}

export interface LocalBranch {
  type: "local"
  name: string
  gitName: string
}

export interface RemoteBranch {
  type: "remote"
  name: string
  gitName: string
  remoteName: string
}

export type Branch = LocalBranch | RemoteBranch

export interface RepoStatus {
  localBranches: LocalBranchStatus[]
  hasUncommitedChanges: boolean
  hasUntrackedFiles: boolean
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

export type RepoType = "standard" | "bare"

export interface DirInfo {
  repoType: RepoType
  repoRoot: string
  worktreeRoot?: string

  isRepoRoot: boolean
  isWorktreeRoot: boolean
}

export interface Worktree {
  fullPath: string
  // TODO: relativePath, localBranch
}
