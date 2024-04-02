import { Repo } from "./Repo.ts"
import { parseBranchLine } from "./helpers.ts"
import { Branch } from "./types.ts"

export class RepoBranch {
  constructor(private repo: Repo) {}

  current(): string {
    return this.repo.git(["rev-parse", "--abbrev-ref", "HEAD"]).stdout
  }

  list(): Branch[] {
    const { stdout } = this.repo.git(["branch", "--all"])
    return stdout
      .split("\n")
      .filter((line) => !/\/HEAD /.test(line)) // ignore HEAD
      .map(parseBranchLine)
  }

  deleteLocalBranch(
    name: string,
    { force = false }: { force?: boolean } = {},
  ): void {
    this.repo.git(["branch", force ? "-D" : "-d", name])
  }

  deleteRemoteBranch(
    { name, remoteName }: { name: string; remoteName: string },
  ): void {
    this.repo.git(["push", "--delete", remoteName, name])
  }
}
