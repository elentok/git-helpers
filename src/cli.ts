import { Command } from "npm:commander@11.1.0"
import { destroy } from "./commands/destroy.ts"
import { status } from "./commands/status.ts"

const program = new Command()
program.command("destroy").description("Destroys a local and remote branch").action(destroy)
program.command("status [branch]").description("Shows the branch sync status").action(status)

program.parse()
