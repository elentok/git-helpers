import { run } from "./run.ts"
import { Repo } from "./types.ts"

export interface InitRepoOptions {
  bare?: boolean
}

export function init(
  root: string,
  { bare = false }: InitRepoOptions = {},
): Repo {
  const args = ["init"]
  if (bare) {
    args.push("--bare")
  }

  run(root, args)

  return { root, isBare: bare }
}
