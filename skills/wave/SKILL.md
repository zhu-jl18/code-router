---
name: wave
description: "Execution-phase parallel acceleration strategy. Use when the host agent already understands the task and needs to execute work efficiently via iterative rounds of flat-parallel dispatch. Each wave: decompose remaining work into independent tasks → dispatch via code-dispatcher --parallel (no dependencies) → review results → plan next wave. Depends on code-dispatcher or code-dispatcher-flash skill."
---

# /wave — Iterative Flat-Parallel Execution

## Rules

1. Tasks within a wave are flat parallel — no `dependencies` field
2. Review results after each wave before dispatching the next
3. Host agent may handle some work directly alongside dispatched tasks

## Execution Loop

```
repeat:
  1. Decompose remaining work into independent tasks
  2. Dispatch wave via --parallel
  3. Review results
  4. If done → break; else → plan next wave from results
```

### Dispatch

```bash
code-dispatcher --parallel --backend {backend} <<'EOF'
---TASK---
id: w{N}-task-1
backend: {backend}
workdir: .
---CONTENT---
[task description, file scope, test command]

---TASK---
id: w{N}-task-2
backend: {backend}
workdir: .
---CONTENT---
[task description, file scope, test command]
EOF
```

Task IDs: `w{wave}-task-{n}` (e.g. `w1-task-1`, `w2-task-3`).

### Review

After each wave:
- **Pass** → plan next wave from what was accomplished
- **Fail** → retry failed tasks (max 2 retries), then move on
- **Conflict** → add merge/fix task to next wave
- **All done** → exit loop, run final verification

### Decomposition Principle

A task belongs in **this wave** only if it needs nothing from other tasks in the same wave. Anything that depends on this wave's output goes to the next wave.
