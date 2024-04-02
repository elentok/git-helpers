import * as git from "./git/index.ts"
import * as path from "https://deno.land/std@0.214.0/path/mod.ts"

export function createDummyRepo(opts?: git.InitRepoOptions): string {
  const dir = Deno.makeTempDirSync()
  git.init(dir, opts)
  return dir
}

export function createDummyCommit(
  repoRoot: string,
  files: Record<string, string> = { "file1.txt": "Hello World" },
) {
  Object.entries(files).forEach(([filename, content]) => {
    const fullname = path.join(repoRoot, filename)
    Deno.writeTextFileSync(fullname, content)
    git.run(repoRoot, ["add", filename])
  })
}
