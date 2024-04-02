import { fzf, FzfFlags } from "./fzf.ts"
import { isPresent } from "./helpers.ts"
import * as git from "./git/index.ts"
import { Branch, Repo } from "./git/types.ts"

export async function pickRemote(
  repo: string | Repo,
  flags?: FzfFlags,
): Promise<string> {
  return (await fzf({
    items: git.remote.list(repo),
    prompt: "Pick remote: ",
    ...flags,
  }))[0]
}

export async function pickBranch(
  repo: Repo,
  flags?: FzfFlags,
): Promise<Branch[]> {
  const branches = git.branch.list(repo)

  const gitNames = await fzf({
    items: branches.map((b) => b.gitName),
    prompt: flags?.allowMultiple ? "Pick branches: " : "Pick branch: ",
    ...flags,
  })

  return gitNames.map((gitName) => branches.find((b) => b.gitName === gitName))
    .filter(isPresent)
}
