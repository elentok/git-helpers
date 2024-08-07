import { parseBranchLine } from "./parseBranchLine.ts"
import { run } from "./run.ts"
import { Branch, Repo } from "./types.ts"

export function create(repo: string | Repo, branchName: string): void {
  run(repo, ["checkout", "-b", branchName])
}

export function current(repo: Repo): string {
  return run(repo, ["rev-parse", "--abbrev-ref", "HEAD"]).stdout
}

export function list(repo: string | Repo): Branch[] {
  const { stdout } = run(repo, ["branch", "--all"])
  return stdout
    .split("\n")
    .filter((line) => !/\/HEAD /.test(line)) // ignore HEAD
    .map(parseBranchLine)
}

export function deleteLocalBranch(
  repo: string | Repo,
  name: string,
  { force = false }: { force?: boolean } = {},
): void {
  run(repo, ["branch", force ? "-D" : "-d", name])
}

export function deleteRemoteBranch(
  repo: string | Repo,
  { name, remoteName }: { name: string; remoteName: string },
): void {
  run(repo, ["push", "--delete", remoteName, name])
}
