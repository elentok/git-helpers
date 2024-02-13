import { git } from "./git.ts"
import { parseBranchLine } from "./helpers.ts"
import { identifyDir } from "./identify-dir.ts"
import { ShellResult } from "./shell.ts"
import { ShellOptions } from "./shell.ts"
import { Branch } from "./types.ts"

export interface CreateRepoOptions {
  bare?: boolean
}

export class Repo {
  private constructor(
    public readonly root: string,
    public readonly isBare: boolean,
  ) {}

  static fromPath(path: string): Repo | undefined {
    const dirInfo = identifyDir(path)
    if (dirInfo == null) return

    const { repoRoot, repoType } = dirInfo

    return new Repo(repoRoot, repoType === "bare")
  }

  static init(
    dir: string,
    { bare = false }: CreateRepoOptions = {},
  ): Repo {
    const args = ["init"]
    if (bare) {
      args.push("--bare")
    }

    git(dir, args)

    return new Repo(dir, bare)
  }

  git(args: string[], options?: ShellOptions): ShellResult {
    return git(this.root, args, options)
  }

  remotes(): string[] {
    return this.git(["remote"]).stdout.split("\n")
  }

  remoteUpdate(): void {
    console.info("Updating remotes...")
    this.git(["remote", "update"])
  }

  hash(ref: string): string {
    return this.git(["log", "-1", "--pretty=%H", ref]).stdout
  }

  currentBranch(): string {
    return this.git(["rev-parse", "--abbrev-ref", "HEAD"]).stdout
  }

  revCount(
    fromRef: string,
    toRef: string,
  ): number {
    const output =
      this.git(["rev-list", "--count", `${fromRef}..${toRef}`]).stdout
    const count = Number(output)

    if (isNaN(count)) {
      throw new Error(`Invalid rev-list count '${output}'`)
    }

    return count
  }

  hasUncommitedChanges(): boolean {
    const output = this.git(["status", "--porcelain=v1"]).stdout
    if (output.length === 0) return false

    const lines = output.split("\n").filter((l) =>
      l.length > 0 && !l.startsWith("?? ")
    )
    return lines.length > 0
  }

  hasUntrackedFiles(): boolean {
    const output = this.git(["status", "--porcelain=v1"]).stdout
    if (output.length === 0) return false

    const lines = output.split("\n")
    return lines.find((l) => l.startsWith("?? ")) != null
  }

  branches(): Branch[] {
    const { stdout } = this.git(["branch", "--all"])
    return stdout
      .split("\n")
      .filter((line) => !/\/HEAD /.test(line)) // ignore HEAD
      .map(parseBranchLine)
  }
}

/**
 * Returns the repository in the given directory (searches up until it finds
 * the root directory).
 *
 * If no repo found shows and error message and exits the process
 */
export function findRepoOrExit(dir: string): Repo {
  const repo = Repo.fromPath(dir)
  if (repo == null) {
    console.error(`Error: No git repo found at '${dir}'`)
    Deno.exit(1)
  }

  return repo
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
