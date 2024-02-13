// import { shell } from "./shell.ts"

import { Branch } from "./types.ts"

export const CHECKMARK = "✔"
export const ERROR = "✘"

export function isPresent<T>(value: T | null | undefined): value is T {
  return value != null
}

export function isDirectory(dir: string): boolean {
  try {
    return Deno.statSync(dir).isDirectory
  } catch (e) {
    console.log("[elentok] [helpers.ts] isDirectory exception", e)
    return false
  }
}

export function parseBranchLine(line: string): Branch {
  line = line.replace(/^\s*\*/g, "").trim()

  if (line.match(/^remotes\//)) {
    const [_, remoteName, name] = line.split("/", 3)
    return {
      type: "remote",
      name,
      remoteName,
      gitName: `${remoteName}/${name}`,
    }
  } else {
    return {
      type: "local",
      name: line,
      gitName: line,
    }
  }
}

// export function hasCommand(command: string): boolean {
//   try {
//     shell("which", {args: command})
//   } catch(e) {
//
//   }
// }
