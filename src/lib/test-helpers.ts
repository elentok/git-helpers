import { git } from "./git.ts"
import { createRepo, CreateRepoOptions } from "./repo.ts"
import * as path from "https://deno.land/std@0.214.0/path/mod.ts"

export function createDummyRepo(opts?: CreateRepoOptions): string {
  const dir = Deno.makeTempDirSync()
  createRepo(dir, opts)
  return dir
}

export function createDummyCommit(
  repoRoot: string,
  files: Record<string, string> = { "file1.txt": "Hello World" },
) {
  Object.entries(files).forEach(([filename, content]) => {
    const fullname = path.join(repoRoot, filename)
    Deno.writeTextFileSync(fullname, content)
    git(repoRoot, ["add", filename])
  })
}
