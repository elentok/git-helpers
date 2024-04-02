import { run } from "./run.ts"

export function revParseString(
  repoRoot: string,
  what: "show-toplevel" | "git-dir",
): string | undefined {
  const result = run(repoRoot, ["rev-parse", `--${what}`], {
    throwError: false,
  })
  if (!result.success) return
  return result.stdout
}

export function revParseBoolean(
  repoRoot: string,
  what: "is-bare-repository" | "is-inside-work-tree",
): boolean {
  const result = run(repoRoot, ["rev-parse", `--${what}`], {
    throwError: false,
  })
  return result.success && result.stdout === "true"
}
