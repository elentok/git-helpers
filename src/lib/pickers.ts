import { gitBranches } from "./branch.ts"
import { FzfFlags, fzf } from "./fzf.ts"
import { gitRemotes } from "./git.ts"
import { isPresent } from "./helpers.ts"
import { Branch, Repo } from "./types.ts"

export async function pickRemote(repo: Repo, flags?: FzfFlags): Promise<string> {
  return (await fzf({ items: gitRemotes(repo), prompt: "Pick remote: ", ...flags }))[0]
}

export async function pickBranch(repo: Repo, flags?: FzfFlags): Promise<Branch[]> {
  const branches = gitBranches(repo)

  const gitNames = await fzf({
    items: branches.map((b) => b.gitName),
    prompt: flags?.allowMultiple ? "Pick branches: " : "Pick branch: ",
    ...flags,
  })

  return gitNames.map((gitName) => branches.find((b) => b.gitName === gitName)).filter(isPresent)
}
