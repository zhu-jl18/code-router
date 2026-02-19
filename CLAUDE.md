# code-dispatcher Project Instructions

## Project Summary
code-dispatcher is a Go CLI that receives coding tasks, selects a backend (`codex`/`claude`/`gemini`), builds invocation parameters, dispatches execution, and collects results.

## Tech Stack
- Go: main dispatcher program
- Python: install/uninstall scripts
- Bash: build and distribution script

## Repository Structure
- `code-dispatcher/`: Go source (main package and backend dispatch logic)
- `skills/`: Claude Code skills (`dev`, `wave`, `code-dispatcher`, `code-dispatcher-flash`, `code-council`, `github-issue-pr-flow`, `pr-review-reply`)
- `docs/`: documentation (`runtime-config.md`)
- `memory/`: additional instruction memory (`CLAUDE-add.md`)
- `scripts/`: build scripts (`build-dist.sh`)
- `install.py` / `uninstall.py`: installer and uninstaller
- `dist/`: build outputs (gitignored)

## Development Commands
```bash
cd code-dispatcher && go test ./...
bash scripts/build-dist.sh
python3 install.py
```

## Code Conventions
- All backend execution must go through `code-dispatcher`; do not call `codex`, `claude`, or `gemini` directly.
- Single runtime config source: `~/.code-dispatcher/.env`.
- Prompt injection files: `~/.code-dispatcher/prompts/<backend>-prompt.md`.

## /dev Workflow Contract (memory/CLAUDE-add.md)
When `/dev ...` is triggered (or `code-dispatcher` is explicitly mentioned):
- Claude Code is responsible for intake, context gathering, planning, and verification.
- Editing files and running tests must be executed through the `code-dispatcher` skill.
