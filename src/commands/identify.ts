import { identifyDirectory } from "../lib/repo.ts"

export function identify() {
  const dirInfo = identifyDirectory(Deno.cwd())
  console.log("[david] [identify.ts] identify", dirInfo)

  if (dirInfo.isBare) {
    console.info("Bare repository")
  } else {
    if (dirInfo.isInsideWorktree) {
      console.info("Inside worktree")
    } else {
      console.info("Regular repository")
    }
  }
}
