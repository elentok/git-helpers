import { shell, ShellOptions, ShellResult } from "../shell.ts"
import { Repo } from "./types.ts"

export function run(
  repo: Repo,
  args: string[],
  options?: ShellOptions,
): ShellResult {
  const repoRoot = typeof repo === "string" ? repo : repo.root
  return shell("git", { args, cwd: repoRoot, ...options })
}
