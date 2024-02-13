import { shell, ShellOptions, ShellResult } from "./shell.ts"

export function git(
  repoRoot: string,
  args: string[],
  options?: ShellOptions,
): ShellResult {
  return shell("git", { args, cwd: repoRoot, ...options })
}

export function revParseString(
  repoRoot: string,
  what: "show-toplevel" | "git-dir",
): string | undefined {
  const result = git(repoRoot, ["rev-parse", `--${what}`], {
    throwError: false,
  })
  if (!result.success) return
  return result.stdout
}

export function revParseBoolean(
  repoRoot: string,
  what: "is-bare-repository" | "is-inside-work-tree",
): boolean {
  const result = git(repoRoot, ["rev-parse", `--${what}`], {
    throwError: false,
  })
  return result.success && result.stdout === "true"
}
