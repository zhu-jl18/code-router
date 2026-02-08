---
name: fish-agent-wrapper
description: Execute fish-agent-wrapper for multi-backend AI code tasks. Supports Codex, Claude, Gemini, Ampcode, with file references (@syntax) and structured output.
---

# fish-agent-wrapper Integration

## Overview

Execute fish-agent-wrapper commands with pluggable AI backends(Codex, Claude, Gemini, Ampcode). Supports file references via `@` syntax, parallel task execution with backend selection, and configurable security controls.

## When to Use

When the user explicitly requests a specific backend (Codex, Claude, Gemini, or Ampcode), mentions fish-agent-wrapper, or when a skill or command definition explicitly declares a dependency on this skill.

Applicable scenarios include but are not limited to:
- Complex code analysis requiring deep understanding
- Large-scale refactoring across multiple files
- Automated code generation with backend selection


## Typical Usage for One-Round/New Tasks:

**1) Standard invocation: HEREDOC syntax (recommended)**
```bash
fish-agent-wrapper --backend codex - [working_dir] <<'EOF'
<task content here>
EOF
```

**2) Single-line tasks (no heredoc)**
```bash
fish-agent-wrapper --backend codex "simple task" [working_dir]
fish-agent-wrapper --backend claude "simple task" [working_dir]
fish-agent-wrapper --backend gemini "simple task" [working_dir]
fish-agent-wrapper --backend ampcode "simple task" [working_dir]
```

## Common Parameters

- Command notation and positional order
  - `[]` means optional. Do not type brackets literally.
  - Inline task: `fish-agent-wrapper --backend <backend> "<task>" [working_dir]` — task text is passed in command args; best for short one-line prompts.
  - Stdin task: `fish-agent-wrapper --backend <backend> - [working_dir]` — task text is read from stdin (`<<'EOF'`/pipe); best for multi-line or complex content.
  - These two forms are for new-task commands only.
  - Resume uses its own forms: inline `fish-agent-wrapper --backend <backend> resume <session_id> "<follow-up task>"` and stdin `fish-agent-wrapper --backend <backend> resume <session_id> -`.

- `--backend <backend>` (required)
  - Select backend explicitly: `codex | claude | gemini | ampcode`.
  - Must be present in both new and resume modes.
  - In parallel mode, this is the global default backend.
  - If a task block defines `backend`, it overrides the global default for that task.

- `task` (required)
  - Task description for the backend.
  - Supports inline text or stdin marker `-`.
  - Supports `@file` references.

- `working_dir` (optional)
  - Working directory for new task execution.
  - Omit it to use the current directory.
  - Typical values: `.`, `./subdir`, `/absolute/path`.
  - In resume mode, do not append `working_dir`; resume follows backend session context.

- Output modes (parallel execution only)
  - **Summary (default)**: Structured report with changes, output, verification, and review summary.
  - **Full (`--full-output`)**: Complete task messages. Use only when debugging specific failures.
  - Scope: `--full-output` is valid only with `--parallel`; single-task mode does not support this flag.
  - Backend behavior: mode selection is wrapper-level and works the same for `codex | claude | gemini | ampcode`.

## Return Format:

```
Agent response text here...

---
SESSION_ID: 019a7247-ac9d-71f3-89e2-a823dbd8fd14
```

## Backends Selection Guide

**Note**: This backends selection guide applies only when the user has not explicitly requested a specific backend. If the user specifies a backend, always follow the user's instructions.


Quiklook at the differences between backends:

| Backend | Command | Description | Best For |
|---------|---------|-------------|----------|
| codex | `--backend codex` | OpenAI Codex (default) | Code analysis, complex development, debugging |
| claude | `--backend claude` | Anthropic Claude | Quick fixes, documentation, prompts |
| gemini | `--backend gemini` | Google Gemini | UI/UX prototyping |
| ampcode | `--backend ampcode` | Sourcegraph Amp | Review tasks, debugging |

For detailed guidance:

**Codex**:
- Deep code understanding and complex logic implementation
- Large-scale refactoring with precise dependency tracking
- Algorithm optimization and performance tuning
- Example: "Analyze the call graph of @src/core and refactor the module dependency structure"

**Claude**:
- Quick feature implementation with clear requirements
- Technical documentation, API specs, README generation
- Professional prompt engineering (e.g., product requirements, design specs)
- Example: "Generate a comprehensive README for @package.json with installation, usage, and API docs"

**Gemini**:
- UI component scaffolding and layout prototyping
- Design system implementation with style consistency
- Interactive element generation with accessibility support
- Example: "Create a responsive dashboard layout with sidebar navigation and data visualization cards"

**Ampcode**:
- Fast code reviews and improvement suggestions
- Plan review and feedback for development proposals
- Debugging assistance when Codex fails.
- Example: "Review @.claude/specs/auth/dev-plan.md and give feedback on potential issues and improvements"

A Typical Backend Switching Example:
- Start with Codex for analysis, switch to Claude for documentation, then Gemini for UI implementation. Use Ampcode for supplementary tasks such as plan review and suggestions.
- Use per-task backend selection in parallel mode to optimize for each task's strengths

## Resume Session

All four backends support resume mode: `codex | claude | gemini | ampcode`.

**1) Standard resume (HEREDOC)**
```bash
fish-agent-wrapper --backend codex resume <session_id> - <<'EOF'
<follow-up task>
EOF
```

**2) Single-line resume (no heredoc)**
```bash
fish-agent-wrapper --backend claude resume <session_id> "follow-up task"
```

**3) Parallel resume (supported)**
```bash
fish-agent-wrapper --parallel --backend codex <<'EOF'
---TASK---
id: resume-a
backend: claude
session_id: sid_claude_1
---CONTENT---
follow-up for claude session

---TASK---
id: resume-b
backend: ampcode
session_id: T-amp-1
---CONTENT---
follow-up for ampcode session
EOF
```

In parallel mode, any task that provides `session_id` runs in resume mode.

Resume mode relies on backend session context.
- Do not append `[working_dir]` in resume commands.
- If you need a different directory, start a new session instead of resume.

## Parallel Execution

Parallel mode uses a dependency DAG scheduler.

- `id` defines a unique task node.
- `dependencies: a, b` means this task waits for tasks `a` and `b` to succeed.
- Tasks in the same DAG layer run concurrently; the next layer starts only after the current layer finishes.
- If a dependency fails, dependent tasks are skipped.
- Invalid dependency IDs or dependency cycles fail fast before execution starts.
- `--backend` in parallel mode is a required global fallback; tasks without `backend` use it. Usally set to `codex`.
- `backend` inside a task block overrides the global fallback for that task.

ASCII execution model:
```text
layer 0: task1      taskX
           |          |
layer 1: task2      taskY
             \      /
layer 2:      task3
```

**1) Dependency scheduling (global backend fallback)**
```bash
fish-agent-wrapper --parallel --backend codex <<'EOF'
---TASK---
id: task1
workdir: /path/to/dir
---CONTENT---
analyze code structure
---TASK---
id: task2
dependencies: task1
---CONTENT---
design architecture based on task1 analysis
EOF
```

**2) Per-task backend override (mixed backends)**
```bash
fish-agent-wrapper --parallel --backend codex <<'EOF'
---TASK---
id: task1
---CONTENT---
analyze code structure
---TASK---
id: task2
backend: claude
dependencies: task1
---CONTENT---
design architecture based on task1 analysis
---TASK---
id: task3
backend: gemini
dependencies: task2
---CONTENT---
generate implementation code
EOF
```

**3) Minimal mixed-backend example (annotated)**
```bash
fish-agent-wrapper --parallel --backend codex <<'EOF'
---TASK---
id: prep
# uses global backend codex
---CONTENT---
scan @src and list key modules

---TASK---
id: plan
backend: claude
# overrides global backend for this task
dependencies: prep
---CONTENT---
write implementation plan based on prep output
EOF
```

In parallel mode, output has two styles:

**1) Summary mode (default, no flag)**
```bash
fish-agent-wrapper --parallel --backend codex <<'EOF'
---TASK---
id: t1
---CONTENT---
analyze @src and summarize architecture changes
EOF
```

**2) Full mode (`--full-output`)**, mainly for debugging failures or when full per-task messages are required.
```bash
fish-agent-wrapper --parallel --backend codex --full-output <<'EOF'
---TASK---
id: t1
---CONTENT---
analyze @src and summarize architecture changes
EOF
```

**Concurrency Control**:
Set `FISH_AGENT_WRAPPER_MAX_PARALLEL_WORKERS` to limit concurrent tasks (default: unlimited).

## Environment Variables

- `CODEX_TIMEOUT`: Override timeout in milliseconds (default: 7200000 = 2 hours)
- `FISH_AGENT_WRAPPER_SKIP_PERMISSIONS`: Control Claude CLI permission checks
  - For **Claude** backend: default is **skip permissions** unless explicitly disabled
  - Set `FISH_AGENT_WRAPPER_SKIP_PERMISSIONS=false` to keep Claude permission prompts
- `FISH_AGENT_WRAPPER_MAX_PARALLEL_WORKERS`: Limit concurrent tasks in parallel mode (default: unlimited, recommended: 8)
- `FISH_AGENT_WRAPPER_CLAUDE_DIR`: Override the base Claude config dir (default: `~/.claude`)
- `FISH_AGENT_WRAPPER_AMPCODE_MODE`: Set Ampcode mode (`smart|deep|rush|free`, default: `smart`)

## Invocation Pattern

**Single Task**:
```
Bash tool parameters:
- command: fish-agent-wrapper --backend <backend> - [working_dir] <<'EOF'
  <task content>
  EOF
- timeout: 7200000
- description: <brief description>

Note: `--backend` is required; supported values: `codex | claude | gemini | ampcode`
```

**Parallel Tasks**:
```
Bash tool parameters:
- command: fish-agent-wrapper --parallel --backend <backend> <<'EOF'
  ---TASK---
  id: task_id
  backend: <backend>  # Optional, overrides global
  workdir: /path
  dependencies: dep1, dep2
  ---CONTENT---
  task content
  EOF
- timeout: 7200000
- description: <brief description>

Note: Global --backend is required; per-task backend is optional
```

## Critical Rules

**NEVER kill fish-agent-wrapper processes.** Long-running tasks are normal. Instead:

1. **Check task status via log file**:
   ```bash
   # View real-time output
   tail -f /tmp/claude/<workdir>/tasks/<task_id>.output

   # Check if task is still running
   cat /tmp/claude/<workdir>/tasks/<task_id>.output | tail -50
   ```

2. **Wait with timeout**:
   ```bash
   # Use TaskOutput tool with block=true and timeout
   TaskOutput(task_id="<id>", block=true, timeout=300000)
   ```

3. **Check process without killing**:
   ```bash
   ps aux | grep fish-agent-wrapper | grep -v grep
   ```

**Why:** fish-agent-wrapper tasks often take 2-10 minutes. Killing them wastes API costs and loses progress.

## Security Best Practices

- **Claude Backend**: Permission checks enabled by default
  - To skip checks: set `FISH_AGENT_WRAPPER_SKIP_PERMISSIONS=true` or pass `--skip-permissions`
- **Concurrency Limits**: Set `FISH_AGENT_WRAPPER_MAX_PARALLEL_WORKERS` in production to prevent resource exhaustion
- **Automation Context**: This wrapper is designed for AI-driven automation where permission prompts would block execution

## Recent Updates

- Multi-backend support for all modes (workdir, resume, parallel)
- Security controls with configurable permission checks
- Concurrency limits with worker pool and fail-fast cancellation
- Ampcode backend support for new/resume/parallel execution
