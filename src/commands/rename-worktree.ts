import { fzf } from "../lib/fzf.ts"
import * as git from "../lib/git/index.ts"
import * as path from "jsr:@std/path"

export async function renameWorktree() {
  const repo = git.findRepoOrExit(Deno.cwd())

  const worktrees = git.worktree.list(repo)
  const selectedItems = await fzf({
    items: worktrees.map((wt) => wt.name),
  })

  if (selectedItems.length !== 1) {
    console.info("User aborted.")
    return
  }

  const fromName = selectedItems[0]
  console.info(`Rename ${fromName}`)
  const toName = prompt("   to:")
  if (toName == null) {
    console.info("User aborted.")
    return
  }

  console.info("- Calling git worktree move")
  git.worktree.move(repo, fromName, toName)

  console.info("- Rename worktrees/ dir")
  const fromFullName = path.join(repo.root, "worktrees", fromName)
  const toFullName = path.join(repo.root, "worktrees", toName)
  Deno.renameSync(fromFullName, toFullName)

  console.info("- Updating gitdir file")
  const gitDirFile = path.join(toFullName, "gitdir")
  const gitDir = Deno.readTextFileSync(gitDirFile).replace(
    `/${fromName}/.git`,
    `/${toName}/.git`,
  )
  Deno.writeTextFileSync(gitDirFile, gitDir)

  console.info("- Updating .git file")
  const dotGitFile = path.join(repo.root, toName, ".git")
  const dotGit = Deno.readTextFileSync(dotGitFile).replace(
    fromFullName,
    toFullName,
  )
  Deno.writeTextFileSync(dotGitFile, dotGit)

  console.info(`- Creating branch '${toName}'`)
  git.branch.create(path.join(git.root(repo), toName), toName)

  console.info("done")
}
