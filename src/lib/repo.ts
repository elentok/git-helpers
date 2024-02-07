import { git, revParseBoolean, revParseString } from "./git.ts"
import { isDirectory } from "./helpers.ts"
import { shell } from "./shell.ts"
import { DirInfo, Repo } from "./types.ts"
import * as path from "https://deno.land/std@0.214.0/path/mod.ts"

/**
 * Returns the repository in the given directory (searches up until it finds
 * the root directory).
 */
export function findRepo(dir: string): Repo | undefined {
  const { success, stdout: root } = shell("git", {
    args: ["rev-parse", "--show-toplevel"],
    cwd: dir,
    throwError: false,
  })
  if (success && isDirectory(root)) {
    return { root }
  }
}

/**
 * Returns the repository in the given directory (searches up until it finds
 * the root directory).
 *
 * If no repo found shows and error message and exits the process
 */
export function findRepoOrExit(dir: string): Repo {
  const repo = findRepo(dir)
  if (repo == null) {
    console.error(`Error: No git repo found at '${dir}'`)
    Deno.exit(1)
  }

  return repo
}

export interface CreateRepoOptions {
  bare?: boolean
}

export function createRepo(
  dir: string,
  { bare = false }: CreateRepoOptions = {},
): void {
  const args = ["init"]
  if (bare) {
    args.push("--bare")
  }

  git(dir, args)
}

/**
 * There are 2 types of repositories:
 *
 * 1. Standard repository
 * 2. Bare repository
 *
 * A bare repository can contain worktrees
 */
export function identifyDirectory(dir: string): DirInfo | undefined {
  const context = getDirContext(dir)
  if (context == null) return

  if (context.isInsideWorktree) {
    return identifyWorktree(context)
  }

  const repoRoot = context.gitDir == "." ? dir : context.gitDir
  return {
    repoType: "bare",
    repoRoot,
    isRepoRoot: repoRoot === dir,
    isWorktreeRoot: false,
  }
}

function identifyWorktree(context: DirContext): DirInfo | undefined {
  const { dir, gitDir, topLevel } = context
  if (topLevel == null) {
    console.warn(
      `Inside worktree, but "git rev-parse --show-toplevel" in ${dir} failed`,
    )
    return
  }

  const gitDirName = path.basename(gitDir)

  // Inside regular repo
  if (gitDirName === ".git") {
    const repoRoot = topLevel
    return {
      repoType: "standard",
      repoRoot,
      isRepoRoot: repoRoot === dir,
      isWorktreeRoot: repoRoot === dir,
    }
  }

  // Inside worktree within a bare repo
  const worktreeRoot = topLevel
  const repoRoot = findRepoRootFromWorktreeRoot(worktreeRoot)
  if (repoRoot == null) {
    console.warn(`Cannot find repo root for worktree "${worktreeRoot}"`)
    return
  }

  return {
    repoType: "bare",
    repoRoot,
    isRepoRoot: repoRoot === dir,
    isWorktreeRoot: worktreeRoot === dir,
    worktreeRoot,
  }
}

function findRepoRootFromWorktreeRoot(
  worktreeRoot: string,
): string | undefined {
  const parentDir = path.dirname(worktreeRoot)
  const parentGitDir = revParseString(parentDir, "git-dir")
  if (parentGitDir == null) {
    console.warn(
      `Inside worktree, but "git rev-parse --git-dir" in ${parentDir} failed`,
    )
    return
  }

  return parentGitDir === "." ? parentDir : parentGitDir
}

function getDirContext(dir: string): DirContext | undefined {
  const gitDir = revParseString(dir, "git-dir")
  if (gitDir == null) return

  return {
    dir,
    gitDir,
    topLevel: revParseString(dir, "show-toplevel"),
    isInsideWorktree: revParseBoolean(dir, "is-inside-work-tree"),
    isBare: revParseBoolean(dir, "is-bare-repository"),
  }
}

interface DirContext {
  dir: string
  gitDir: string
  topLevel?: string
  isInsideWorktree: boolean
  isBare: boolean
}

// import * as shell from "shelljs"
// import { IPair, LocalBranch, parseBranchLine, RemoteBranch } from "./branch"
// import { Remote } from "./remote"
// import { IRepo } from "./types"
//
// export class Repo implements IRepo {
//   public remotes: Remote[]
//   public localBranchesByName: { [key: string]: LocalBranch } = {}
//   public remoteBranchesByName: { [key: string]: RemoteBranch[] } = {}
//
//   constructor(public root: string) {
//     this.remotes = this.loadRemotes()
//     this.loadBranches()
//   }
//
//   public localBranches(): LocalBranch[] {
//     return Object.values(this.localBranchesByName)
//   }
//
//   public remoteBranches(): RemoteBranch[] {
//     return ([] as RemoteBranch[]).concat(...Object.values(this.remoteBranchesByName))
//   }
//
//   public findLocalBranchByName(name: string): LocalBranch {
//     if (this.localBranchesByName == null) this.loadBranches()
//     return this.localBranchesByName[name]
//   }
//
//   public findRemoteBranchesByName(name: string): RemoteBranch[] {
//     if (this.remoteBranchesByName == null) this.loadBranches()
//     return this.remoteBranchesByName[name]
//   }
//
//   public fetchRemotes(): void {
//     this.remotes
//       .filter((r) => r.name !== "review")
//       .forEach((r) => {
//         r.fetch()
//         r.prune()
//       })
//   }
//
//   public unsyncedBranches(): IPair[] {
//     const pairs: IPair[] = []
//
//     this.localBranches().forEach((local) => {
//       local.remoteBranches.forEach((remote) => {
//         if (local.hash() !== remote.hash()) {
//           pairs.push({ local, remote })
//         }
//       })
//     })
//
//     return pairs
//   }
//
//   public git(command: string, options: shell.ExecOptions = {}): string {
//     shell.cd(this.root)
//     const result = shell.exec(`git ${command}`, options) as shell.ExecOutputReturnValue
//
//     if (result.code !== 0) {
//       throw new Error(`Git command returns status ${result.code}:\n${result.stderr.toString()}`)
//     }
//
//     return result.stdout.toString().trim()
//   }
//
//   private loadRemotes(): Remote[] {
//     return this.git("remote", { silent: true })
//       .split("\n")
//       .map((name) => new Remote(this, name))
//   }
//
//   private loadBranches(): void {
//     this.git("branch --all", { silent: true })
//       .split("\n")
//       .forEach((line) => {
//         const branch = parseBranchLine(line, this)
//         if (branch.name === "HEAD") return
//
//         if (branch instanceof RemoteBranch) {
//           if (this.remoteBranchesByName[branch.name] == null) {
//             this.remoteBranchesByName[branch.name] = []
//           }
//           this.remoteBranchesByName[branch.name].push(branch)
//         } else {
//           this.localBranchesByName[branch.name] = branch
//         }
//       })
//
//     this.addRemotesToLocalBranches()
//   }
//
//   private addRemotesToLocalBranches(): void {
//     Object.keys(this.remoteBranchesByName).forEach((name) => {
//       this.remoteBranchesByName[name].forEach((remoteBranch) => {
//         const localBranch = this.localBranchesByName[name]
//         if (localBranch != null) {
//           localBranch.remoteBranches.push(remoteBranch)
//         }
//       })
//     })
//   }
// }
