import { identifyDir } from "./identifyDir.ts"
import { RepoWithDetails } from "./types.ts"

export function findRepo(path: string): RepoWithDetails | undefined {
  const dirInfo = identifyDir(path)
  if (dirInfo == null) return

  return dirInfo.repo
}

export function findRepoOrThrow(dir: string): RepoWithDetails {
  const repo = findRepo(dir)
  if (repo == null) {
    throw new Error(`No git repo found at '${dir}'`)
  }

  return repo
}

/**
 * Returns the repository in the given directory (searches up until it finds
 * the root directory).
 *
 * If no repo found shows and error message and exits the process
 */
export function findRepoOrExit(dir: string): RepoWithDetails {
  const repo = findRepo(dir)
  if (repo == null) {
    console.error(`Error: No git repo found at '${dir}'`)
    Deno.exit(1)
  }

  return repo
}
