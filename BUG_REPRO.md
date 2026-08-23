# Bug Reproduction

This record is reproduced against the published red branch and verified against the green branch.

- Bug: alert-orchestration-context-cancellation-boundary-002
- Red verification: `go test ./internal/errchain -run '^TestR001Status$' -count=1`
- Green verification: the same acceptance commands return exit code 0 after the production fix.
- Evidence: `_evidence/red_result.json`, `_evidence/verify_green.jsonl`, and `_evidence/verify_result.json`.
