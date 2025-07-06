import { identifyDir } from "../lib/git/identifyDir.ts"
import { log } from "../lib/git/log.ts"

export function changelog(from: string, to: string) {
  const dirInfo = identifyDir(Deno.cwd())
  if (dirInfo == null) {
    console.info("Not a repository")
    Deno.exit(1)
  }

  const commits = log(dirInfo.repo, from, to)
  console.log("[bazinga] [changelog.ts] L12", commits)
}
