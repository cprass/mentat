# TODO

- Add CLI command to add a new deck
- Add CLI command and interactive-mode keybind to add new cards to a deck
- Surface errors to the user when SyncPush fails (currently silent)
- Warn the user clearly when ff-only pull fails instead of silently erroring
- Guard against broad `git add .` committing unexpected files (e.g. .gitignore or filtering)
- auto-sync only in interactive mode, not for the other commands
- add CLI command to run auto-sync
- auto-sync: checkout branch before committing or pushing anything

## Ideas

- Consider moving config.yaml into the vault so it syncs with the rest of the data. Unclear what the implications are (e.g. machine-specific settings, bootstrap problem of needing config to find the vault)
