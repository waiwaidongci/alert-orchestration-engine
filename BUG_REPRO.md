# Bug Reproduction

This record is reproduced against the published red branch and verified against the green branch.

- Bug: alert-orchestration-shared-counter-locking-008
- Red verification: the four `go test -race ./internal/sharedrace -run '^TestR008...' -count=1` acceptance commands report races or fail under concurrent access.
- Green verification: `go test -race ./internal/sharedrace -run '^TestR008' -count=1` and `go test ./...` return exit code 0 after the production fixes.
- Evidence: `_evidence/red_result.json`, `_evidence/verify_green.jsonl`, and `_evidence/verify_result.json`.
