---
name: github-issue-pr-flow
description: "Autonomous GitHub issue-to-PR delivery workflow with bot-review triage and squash-merge closure. Use when handling end-to-end issue-driven implementation, PR lifecycle automation, or bot review rebut/fix loops."
---

# GitHub Issue → PR Delivery Flow

## Operating Principles
- Keep full autonomy for decomposition, branch planning, triage, and merge execution.
- Keep traceability: every PR body must include `Closes #<primary_issue>`.
- Keep scope discipline: one PR should close one primary issue chain.
- Use `code-dispatcher --backend codex` only as backup for `critical`/`high` findings or hard debugging.

## Context Discovery
- Detect repository and default base branch from `gh repo view`.
- Do not ask user for repo/base unless `gh` metadata is unavailable.
- Load templates from `references/issue-templates.md` and review reply patterns from `references/review-playbook.md` when needed.

## Autonomous Workflow
Execute phases in order with no user confirmation gates unless escalation rules are hit.

### Phase 0: Sync Base Branch
- Fetch and fast-forward base before creating or updating working branches.
- If `git pull --ff-only` fails:
1. Stop automatic progression for this branch.
2. Report the exact failure and conflicting refs to user.
3. Ask user to choose between: manual conflict handling, creating a fresh branch from remote base, or aborting this run.
4. Never perform automatic hard reset or force push.

### Phase 1: Decompose Into Issues
- Decompose by deliverables, not by file counts.
- Prefer one issue = one user-verifiable outcome with acceptance criteria and test evidence.
- Split work when merge/review risk is high or when tasks can land independently.
- Create epic only when multiple child issues need a parent umbrella.

### Phase 2: Implement With Branch Discipline
- Create branch names as `feat/issue-<number>-<slug>` or `fix/issue-<number>-<slug>`.
- Keep one active PR branch per issue chain.
- Keep commits scoped to issue intent; avoid mixed-purpose commits.
- Follow `references/branch-history.md` for squash-merge aftermath and new-PR history hygiene.

### Phase 3: Open PR and Link Issues
- Ensure PR body contains `Closes #<primary>` and optional `Relates #<secondary>`.
- Keep PR summary testable and concise.

### Phase 4: Collect Reviews and CI Signals
- Wait a reasonable bot review window, then collect latest review and check-run states.
- Do not block forever if bots stay silent; CI and maintainers remain the gate.

### Phase 5: Triage Review Findings
- Decide each finding independently: fix, rebut with evidence, or align with repo convention.
- Reply under each review comment and resolve each thread after handling.
- Re-run CI after each push; keep looping until stable or escalation is required.
- Limit autonomous rebut/fix loops to 2 rounds, then escalate if blocking findings remain.

### Phase 6: Squash Merge and Closure
- Merge using squash and delete remote branch.
- Verify linked issues are auto-closed; close residual issues manually if needed.
- Apply post-merge local cleanup and next-PR branch reset rules from `references/branch-history.md`.

## Escalation Policy
Escalate to user only when:
1. Base branch sync cannot be fast-forwarded.
2. A destructive action is required (`force-push`, history rewrite on shared branch, closing unrelated issues).
3. Two triage rounds still leave blocking review/CI failures.
4. Fix requires unavailable credentials or infrastructure access.

Everything else should proceed autonomously.

## Failure Handling
- Bot reviews absent: proceed with CI and maintainer review signals.
- Bot feedback conflicts: prioritize reproducible failures and test evidence.
- Scope creep appears mid-flight: split to new issue(s) and keep current PR focused.
- `gh` command fails: retry once; if still failing, report the blocking error and stop.
