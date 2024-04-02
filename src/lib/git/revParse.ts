import { run } from "./run.ts"
import { Repo } from "./types.ts"

export function revParseString(
  repoRoot: string | Repo,
  what: "show-toplevel" | "git-dir",
): string | undefined {
  const result = run(repoRoot, ["rev-parse", `--${what}`], {
    throwError: false,
  })
  if (!result.success) return
  return result.stdout
}

export function revParseBoolean(
  repoRoot: string | Repo,
  what: "is-bare-repository" | "is-inside-work-tree",
): boolean {
  const result = run(repoRoot, ["rev-parse", `--${what}`], {
    throwError: false,
  })
  return result.success && result.stdout === "true"
}
