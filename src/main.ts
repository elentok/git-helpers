import { Repo } from "./lib/types.ts"
import { remotes } from "./lib/remote.ts"
import { branches } from "./lib/branch.ts"

const repo: Repo = { root: Deno.env.get("HOME")! + "/.dotfiles" }

console.info("Remotes:")
for (const remote of remotes(repo)) {
  console.info(`- ${remote}`)
}

console.info("Branches:")
for (const branch of branches(repo)) {
  if (branch.type === "local") {
    console.info(`- LOCAL: ${branch.name} (${branch.gitName})`)
  } else {
    console.info(`- REMOTE (${branch.remoteName}): ${branch.name} (${branch.gitName})`)
  }
}
