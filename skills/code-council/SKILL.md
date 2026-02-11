---
name: code-council
description: Multi-perspective parallel code review using code-router. Runs 2-3 AI reviewers with distinct roles simultaneously, synthesizes findings, then the host agent performs a mandatory final review pass before presenting to the user.
---

# code-council: Multi-Perspective Code Review

## Overview

Leverage code-router's multi-backend parallel execution to run 2-3 AI reviewers simultaneously, each with a distinct review lens. After code-router synthesizes their findings, the **host agent performs a mandatory final review pass** — validating, challenging, and enriching the report with project-wide context before presenting to the user.

## When to Use

- User explicitly requests code review, architecture audit, or quality check
- PR review or pre-merge quality gate
- Module-level or file-level review
- A skill or command definition declares a dependency on this skill

## Review Roles

### Always Active

1. **Logic Reviewer** (preferred backend: `codex`)
   - Design patterns, SOLID violations, DRY
   - Performance bottlenecks, algorithmic complexity
   - Race conditions, concurrency issues
   - Dependency coupling

2. **Spec Reviewer** (preferred backend: `claude`)
   - Edge cases and boundary conditions
   - Error handling completeness
   - Type safety, null safety
   - Naming consistency, API contract clarity
   - Missing or misleading documentation

### Conditionally Active

3. **UX Reviewer** (preferred backend: `gemini`)
   - **Activation rule**: target files contain UI indicators (see detection criteria below)
   - Accessibility (a11y) compliance
   - Responsive design, breakpoint handling
   - UI state completeness (loading, error, empty states)
   - Visual consistency with design system

   **When NOT activated**: only 2 reviewers run. Do NOT force UX review on pure backend/infra/CLI code.

## Workflow

### Step 1: Determine Review Scope

Resolve target files from user input:

- **Explicit references**: user provides file/directory paths → use `@` syntax
- **Git diff**: user mentions PR, diff, or recent changes → capture diff output, pass as context
- **Module**: user names a package or directory → glob relevant source files

If scope is ambiguous, ask the user to clarify.

### Step 2: UI Detection

Scan target file paths and (when practical) content for UI indicators:

```
Match ANY → has_ui = true:
  File extensions : .tsx .jsx .vue .svelte .css .scss .less
  Content patterns: className= | styled\. | @apply | css-modules import | tailwind utility classes
```

Store result as `has_ui`.

### Step 3: Build & Execute Parallel Review

Construct a single `code-router --parallel` invocation.

**Template (3 reviewers — has_ui = true)**:

```bash
code-router --parallel --backend codex <<'EOF'
---TASK---
id: logic
backend: codex
---CONTENT---
You are "The Logic Reviewer."

Review the following code with these priorities:
1. Design patterns and SOLID principle violations
2. Performance bottlenecks and algorithmic complexity
3. Race conditions and concurrency issues
4. Unnecessary coupling and dependency problems

Target: [@ file references]

Output format — use EXACTLY this structure:
## Logic Review Findings
For each finding:
- **[CRITICAL|WARNING|INFO]**: one-line summary
  - Location: file:line or file:function
  - Detail: what is wrong and why it matters
  - Suggestion: concrete fix or improvement

---TASK---
id: spec
backend: claude
---CONTENT---
You are "The Spec Reviewer."

Review the following code with these priorities:
1. Edge cases and boundary conditions not handled
2. Error handling gaps (missing catch, unhandled rejections, unchecked returns, etc.)
3. Type safety issues and potential null/undefined access
4. Naming inconsistencies and API contract clarity
5. Missing, outdated, or misleading documentation

Target: [@ file references]

Output format — use EXACTLY this structure:
## Spec Review Findings
For each finding:
- **[CRITICAL|WARNING|INFO]**: one-line summary
  - Location: file:line or file:function
  - Detail: what is wrong and why it matters
  - Suggestion: concrete fix or improvement

---TASK---
id: ux
backend: gemini
---CONTENT---
You are "The UX Reviewer."

Review the following code with these priorities:
1. Accessibility: missing ARIA labels, keyboard navigation, screen reader support
2. Responsive design: breakpoint handling, mobile-first issues
3. UI state completeness: loading, error, empty, success states all handled
4. Visual consistency with design system conventions

Target: [@ file references]

Output format — use EXACTLY this structure:
## UX Review Findings
For each finding:
- **[CRITICAL|WARNING|INFO]**: one-line summary
  - Location: file:line or file:function
  - Detail: what is wrong and why it matters
  - Suggestion: concrete fix or improvement

---TASK---
id: synthesis
dependencies: logic, spec, ux
backend: claude
---CONTENT---
You are the Review Synthesizer. Merge the review reports from logic, spec, and ux into one unified report.

Rules:
1. Deduplicate: if multiple reviewers flag the same issue, merge into one entry and note which reviewers agree
2. Resolve conflicts: if reviewers disagree, present both perspectives fairly
3. Rank by severity: CRITICAL first, then WARNING, then INFO
4. Group findings by file, not by reviewer
5. Prepend a summary: total counts by severity, top 3 most impactful issues

Output: a single markdown report titled "## Code Council — Synthesized Review"
EOF
```

**Template (2 reviewers — has_ui = false)**:

Same as above but **omit the `ux` task entirely** and change the synthesis dependencies to `logic, spec`.

```bash
code-router --parallel --backend codex <<'EOF'
---TASK---
id: logic
backend: codex
---CONTENT---
[same as above]

---TASK---
id: spec
backend: claude
---CONTENT---
[same as above]

---TASK---
id: synthesis
dependencies: logic, spec
backend: claude
---CONTENT---
You are the Review Synthesizer. Merge the review reports from logic and spec into one unified report.

Rules:
1. Deduplicate: if multiple reviewers flag the same issue, merge into one entry and note which reviewers agree
2. Resolve conflicts: if reviewers disagree, present both perspectives fairly
3. Rank by severity: CRITICAL first, then WARNING, then INFO
4. Group findings by file, not by reviewer
5. Prepend a summary: total counts by severity, top 3 most impactful issues

Output: a single markdown report titled "## Code Council — Synthesized Review"
EOF
```

### Step 4: Host Agent Final Review (MANDATORY — DO NOT SKIP)

After code-router returns, the host agent MUST perform its own review pass. This is the key differentiator of code-council.

**Why this step exists**:
- code-router tasks run in isolation — they lack cross-project context
- The host agent has access to the full codebase, git history, and project conventions
- Individual reviewers may produce false positives or miss project-specific issues

**What the host agent does**:

1. **Read the synthesized report** from code-router output
2. **Validate each finding**:
   - Does this finding actually apply to this codebase's conventions?
   - Is the severity rating appropriate?
   - Mark any false positive as `[DISPUTED]` with a brief reason
3. **Add missed issues**: Use the host agent's broader context to catch things the council missed (e.g., known tech debt, project-specific anti-patterns, cross-module implications)
4. **Produce the final output**: Present findings to the user, organized by priority:
   - Top 3-5 actionable items highlighted at the top
   - Full findings list grouped by file
   - Any disputed items noted with reasoning

### Step 5: Offer Follow-up Actions

After presenting, offer concrete next steps:
- **Fix critical issues**: generate fix tasks and run them through code-router
- **Review more files**: restart the workflow on another target
- **Save report**: write to a file if the user wants a record

## Backend Routing

Default preferred mapping:
- `logic` → `codex`
- `spec` → `claude`
- `ux` → `gemini`
- `synthesis` → `claude`

Fallback rules (same as code-router convention):
- If preferred backend is unavailable, fall back by priority: `codex` → `claude` → `gemini`
- If user says "use only X": all reviewers use that backend — **different prompts still provide distinct perspectives**

Note: Claude is a strong general-purpose reviewer. When only `claude` is available, it handles all roles effectively. Do not treat it as limited to spec review.

## Severity Definitions

- **CRITICAL**: Will cause bugs, security vulnerabilities, data loss, or crashes in production
- **WARNING**: Code smell, maintainability risk, or potential future bug under reasonable conditions
- **INFO**: Style suggestion, minor improvement, or optimization opportunity

## Edge Cases

**Single file, trivial size (<30 lines)**:
- Still run the full council — small code can have critical issues
- Reviewers will naturally produce fewer findings

**No backend explicitly requested by user**:
- Use the default preferred mapping above

**User provides a git diff instead of files**:
- Capture the diff content, pass it inline in the `---CONTENT---` sections instead of `@` file references
- Adjust reviewer prompts to focus on "changes" rather than "code"

**Target has mixed UI and non-UI files**:
- Activate UX reviewer, but scope it to only the UI-relevant files
- Logic and Spec reviewers still review all files

## Critical Rules

1. **NEVER skip Step 4 (host agent final review)** — this is the core value proposition
2. **NEVER force UX review on non-UI code** — respect the detection result
3. **NEVER kill code-router processes** — follow all process management rules from the code-router skill
4. **NEVER modify source code in this skill** — code-council is read-only; fixes go through a separate action

## Example Invocations

**Review specific files**:
```
Review @src/auth/login.ts and @src/auth/session.ts using code-council
```

**Review recent changes (PR-style)**:
```
Run code-council on the uncommitted changes
```

**Audit a module**:
```
code-council audit @src/payments/
```

**Single backend mode**:
```
Run code-council on @src/core/ using only claude
```
