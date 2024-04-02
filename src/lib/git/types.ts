export interface RepoWithDetails {
  root: string
  isBare: boolean
}

export type Repo = string | RepoWithDetails

export type RepoType = "standard" | "bare"

export interface DirInfo {
  repoType: RepoType
  repoRoot: string
  worktreeRoot?: string

  isRepoRoot: boolean
  isWorktreeRoot: boolean
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

export interface Worktree {
  fullPath: string
  name: string
  branchName: string
}
