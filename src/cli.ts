import { Command } from "npm:commander@11.1.0"
import { destroy } from "./commands/destroy.ts"
import { status } from "./commands/status.ts"
import { identify } from "./commands/identify.ts"

const program = new Command()
program.command("destroy").description("Destroys a local and remote branch").action(destroy)
program
  .command("status")
  .option("-q, --quick", "skip updating the remotes")
  .description("Shows the branch sync status")
  .action(status)

program.command("identify").description("Identifies the current directory").action(identify)

program.parse()
