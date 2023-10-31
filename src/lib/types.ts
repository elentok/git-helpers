export interface Repo {
  root: string
}

export interface Remote {
  name: string
  url: string
}

export class GitError extends Error {
  constructor(
    public code: number,
    public command: string,
    public output: string,
  ) {
    super(`Git command '${command}' failed with exitcode ${code}:\n\n${output}`)
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
