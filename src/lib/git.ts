import { GitError, Repo } from "./types.ts"

const decoder = new TextDecoder()

export function git(repo: Repo, args: string[], options?: Deno.CommandOptions): string {
  const command = new Deno.Command("git", { args, cwd: repo.root, ...options })
  const { code, stdout, stderr } = command.outputSync()
  if (code !== 0) {
    const output = decoder.decode(stdout) + "\n" + decoder.decode(stderr)
    throw new GitError(code, ["git", ...args].join(" "), output)
  }

  return decoder.decode(stdout).trim()
}
