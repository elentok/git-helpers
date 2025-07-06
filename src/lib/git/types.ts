export interface Repo {
  root: string
  isBare: boolean
}

export interface DirInfo {
  repo: Repo
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

export interface Commit {
  hash: string
  subject: string
  body: string
}
