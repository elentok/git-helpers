import { Repo } from "./types.ts"

export function findRepo(path: string): Repo | undefined {
  const dirInfo = identifyDir(path)
  if (dirInfo == null) return

  const { repoRoot, repoType } = dirInfo

  return new Repo(repoRoot, repoType === "bare")
}
