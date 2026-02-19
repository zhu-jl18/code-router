# Triage Guide

## Decision Criteria

```text
finding
 ├─ reproducible bug / test failure / policy violation  → fix
 ├─ misunderstanding / false positive (code disproves)  → rebut
 └─ style / nitpick / preference                        → follow repo convention; rebut if conflicts
```

Evidence takes priority over reviewer severity label.

## Fix Reply Template
```
Decision: accepted and fixed

Changes:
- <what changed and why>

Verification:
- <test/command/result>
```

## Rebuttal Template
```
Decision: rebuttal

Reason:
- <concrete evidence: file path / line / test / invariant>

Verification:
- <test/command/result that disproves the concern>
```

## Nitpick / Won't Fix Template
```
Decision: acknowledged, no code change

Reason:
- <one-line explanation aligned with repo convention>
```
