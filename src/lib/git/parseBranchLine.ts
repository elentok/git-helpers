import { Branch } from "./types.ts"

export function parseBranchLine(line: string): Branch {
  line = line.replace(/^\s*[\*\+]/g, "").trim()

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
