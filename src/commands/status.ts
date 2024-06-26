import { CHECKMARK, ERROR } from "../lib/helpers.ts"
import * as git from "../lib/git/index.ts"
import { getRepoStatus, SyncStatus } from "../lib/status.ts"
import chalk from "npm:chalk"

export function status({ quick }: { quick?: boolean } = {}) {
  const dirInfo = git.identifyDir(Deno.cwd())
  if (dirInfo == null) {
    console.error(chalk.red("Error: Not inside a git repository"))
    Deno.exit(1)
  }

  const dirPrettyInfo = [
    dirInfo.isRepoRoot ? "repo root" : null,
    dirInfo.isWorktreeRoot ? "worktree root" : null,
    // dirInfo.isBare ? "bare" : null,
    // dirInfo.isInsideWorktree ? "inside-worktree" : null,
  ]
    .filter((x) => x != null)
    .join(", ")
  console.info(`Repo: ${dirPrettyInfo}`)

  const repo = git.findRepoOrExit(Deno.cwd())
  if (!quick) {
    git.remote.update(repo)
  }
  const status = getRepoStatus(repo)

  if (status.hasUncommitedChanges) {
    console.info(chalk.red(`${ERROR} uncommited changes`))
  }

  if (status.hasUntrackedFiles) {
    console.info(chalk.yellow(`${ERROR} untracked files`))
  }

  for (const localBranch of status.localBranches) {
    if (localBranch.remoteBranches.length === 0) {
      console.info(`- ${localBranch.name} (local only)`)
    } else {
      const symbol = localBranch.isSynced ? CHECKMARK : ERROR

      if (localBranch.remoteBranches.length === 1) {
        const rb = localBranch.remoteBranches[0]
        colorBySyncStatus(
          `${symbol} ${localBranch.gitName} (${rb.status.pretty})`,
          rb.status,
        )
      } else {
        console.info(`${symbol} ${localBranch.gitName}:`)
        for (const rb of localBranch.remoteBranches) {
          const symbol = rb.status.name === "same" ? CHECKMARK : ERROR
          colorBySyncStatus(
            `  ${symbol} ${rb.gitName} (${rb.status.pretty})`,
            rb.status,
          )
        }
      }
    }
  }
}

function colorBySyncStatus(text: string, status: SyncStatus): void {
  switch (status.name) {
    case "behind":
      console.info(chalk.yellow(text))
      break
    case "ahead":
      console.info(chalk.magenta(text))
      break
    case "diverged":
      console.info(chalk.red(text))
      break
    case "same":
      console.info(chalk.green(text))
      break
    case "unclear":
      console.info(chalk.bgRed(text))
      break
  }
}
