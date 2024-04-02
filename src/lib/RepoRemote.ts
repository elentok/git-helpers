import { Repo } from "./Repo.ts"

export class RepoRemote {
  constructor(private repo: Repo) {}

  list(): string[] {
    return this.repo.git(["remote"]).stdout.split("\n")
  }

  update(): void {
    console.info("Updating remotes...")
    this.repo.git(["remote", "update"])
  }
}
