import { git } from "./git.ts"
import { Branch, Repo } from "./types.ts"

export function gitBranches(repo: Repo): Branch[] {
  const { stdout } = git(repo, ["branch", "--all"])
  return stdout
    .split("\n")
    .filter((line) => !/\/HEAD /.test(line)) // ignore HEAD
    .map(_parseBranchLine)
}

export function _parseBranchLine(line: string): Branch {
  line = line.replace(/^\s*\*/g, "").trim()

  if (line.match(/^remotes\//)) {
    const [_, remoteName, name] = line.split("/", 3)
    return {
      type: "remote",
      name,
      remoteName,
      gitName: `${remoteName}/${name}`,
    }
  } else {
    return {
      type: "local",
      name: line,
      gitName: line,
    }
  }
}

// export function parseBranchLine(line: string): LocalBranch | RemoteBranch {
//   line = line.replace(/^\*/g, "").trim()
//   if (line.match(/^remotes\//)) {
//     return new RemoteBranch(repo, line)
//   } else {
//     return new LocalBranch(repo, line)
//   }
// }
//
//
// // import { Branch, Repo } from "./types.ts"
// //
// // abstract class Branch implements IBranch {
// //   public gitName: string
// //   private cachedHash?: string
// //
// //   constructor(
// //     protected repo: IRepo,
// //     public name: string,
// //   ) {
// //     this.gitName = name
// //   }
// //
// //   public hash(): string {
// //     if (this.cachedHash == null) {
// //       this.cachedHash = this.repo.git(`log -1 --pretty=%h ${this.gitName}`, {
// //         silent: true,
// //       })
// //     }
// //     return this.cachedHash
// //   }
// //
// //   public resetHash(): void {
// //     this.cachedHash = undefined
// //   }
// //
// //   public toString(): string {
// //     return this.gitName
// //   }
// //
// //   public isProtected(): boolean {
// //     return ["master", "develop", "protected"].indexOf(this.name) !== null
// //   }
// // }
// //
// // export class LocalBranch extends Branch {
// //   public remoteBranches: RemoteBranch[] = []
// //
// //   public hasRemotes(): boolean {
// //     return this.remoteBranches.length > 0
// //   }
// //
// //   public destroy({ force = false } = {}): string {
// //     const arg = force ? "-D" : "-d"
// //     return this.repo.git(`branch ${arg} ${this.name}`)
// //   }
// // }
// //
// // export class RemoteBranch extends Branch {
// //   public remoteName: string
// //
// //   constructor(repo: IRepo, name: string) {
// //     super(repo, name)
// //
// //     name = name.replace(/ ->.*/g, "")
// //     const parts = name.split("/", 3)
// //     this.remoteName = parts[1]
// //     this.name = parts[2]
// //     this.gitName = `${this.remoteName}/${this.name}`
// //   }
// //
// //   public destroy(): string {
// //     return this.repo.git(`push --delete ${this.remoteName} ${this.name}`)
// //   }
// // }
// //
// // export function parseBranchLine(line: string, repo: IRepo): LocalBranch | RemoteBranch {
// //   line = line.replace(/^\*/g, "").trim()
// //   if (line.match(/^remotes\//)) {
// //     return new RemoteBranch(repo, line)
// //   } else {
// //     return new LocalBranch(repo, line)
// //   }
// // }
// //
// // export interface IPair {
// //   local: LocalBranch
// //   remote: RemoteBranch
// // }
