// import { shell } from "./shell.ts"

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

// export function hasCommand(command: string): boolean {
//   try {
//     shell("which", {args: command})
//   } catch(e) {
//
//   }
// }
