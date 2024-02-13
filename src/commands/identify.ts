import { identifyDir } from "../lib/identify-dir.ts"

export function identify() {
  const dirInfo = identifyDir(Deno.cwd())
  // console.log("[david] [identify.ts] identify", dirInfo)

  if (dirInfo == null) {
    console.info("Not a repository")
  } else {
    console.info(dirInfo)
  }
}
