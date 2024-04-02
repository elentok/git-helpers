import { revParseBoolean, revParseString } from "./revParse.ts"
import { DirInfo } from "./types.ts"
import * as path from "https://deno.land/std@0.214.0/path/mod.ts"

/**
 * There are 2 types of repositories:
 *
 * 1. Standard repository
 * 2. Bare repository
 *
 * A bare repository can contain worktrees
 */
export function identifyDir(dir: string): DirInfo | undefined {
  const context = getDirContext(dir)
  if (context == null) return

  if (context.isInsideWorktree) {
    return identifyWorktree(context)
  }

  const repoRoot = context.gitDir == "." ? dir : context.gitDir
  return {
    repo: {
      isBare: true,
      root: repoRoot,
    },
    isRepoRoot: repoRoot === dir,
    isWorktreeRoot: false,
  }
}

function identifyWorktree(context: DirContext): DirInfo | undefined {
  const { dir, gitDir, topLevel } = context
  if (topLevel == null) {
    console.warn(
      `Inside worktree, but "git rev-parse --show-toplevel" in ${dir} failed`,
    )
    return
  }

  const gitDirName = path.basename(gitDir)

  // Inside regular repo
  if (gitDirName === ".git") {
    const repoRoot = topLevel
    return {
      repo: {
        isBare: false,
        root: repoRoot,
      },
      isRepoRoot: repoRoot === dir,
      isWorktreeRoot: repoRoot === dir,
    }
  }

  // Inside worktree within a bare repo
  const worktreeRoot = topLevel
  const repoRoot = findRepoRootFromWorktreeRoot(worktreeRoot)
  if (repoRoot == null) {
    console.warn(`Cannot find repo root for worktree "${worktreeRoot}"`)
    return
  }

  return {
    repo: {
      isBare: true,
      root: repoRoot,
    },
    isRepoRoot: repoRoot === dir,
    isWorktreeRoot: worktreeRoot === dir,
    worktreeRoot,
  }
}

function findRepoRootFromWorktreeRoot(
  worktreeRoot: string,
): string | undefined {
  const parentDir = path.dirname(worktreeRoot)
  const parentGitDir = revParseString(parentDir, "git-dir")
  if (parentGitDir == null) {
    console.warn(
      `Inside worktree, but "git rev-parse --git-dir" in ${parentDir} failed`,
    )
    return
  }

  return parentGitDir === "." ? parentDir : parentGitDir
}

function getDirContext(dir: string): DirContext | undefined {
  const gitDir = revParseString(dir, "git-dir")
  if (gitDir == null) return

  return {
    dir,
    gitDir,
    topLevel: revParseString(dir, "show-toplevel"),
    isInsideWorktree: revParseBoolean(dir, "is-inside-work-tree"),
    isBare: revParseBoolean(dir, "is-bare-repository"),
  }
}

interface DirContext {
  dir: string
  gitDir: string
  topLevel?: string
  isInsideWorktree: boolean
  isBare: boolean
}
