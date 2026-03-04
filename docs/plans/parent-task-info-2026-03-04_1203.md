# Show Parent Task Info on Subtasks — Completed

## Overview
When viewing a subtask, display parent task info (ID + name) in the detail view.

## Tasks
- [x] Write tests for Task model (Parent field)
- [x] Write tests for FormatTask (parent display)
- [x] Add `Parent` and `ParentTask` fields to Task model
- [x] Update `FormatTask` to display parent info
- [x] Update `getTask` to fetch parent task name
- [x] Run `make check` - all tests pass
- [x] Build for all platforms with `make build-all`

## Progress Log

| DateTime | Task | Status | Notes |
|----------|------|--------|-------|
| 2026-03-04 12:03 | Started implementation | In Progress | TDD approach |
| 2026-03-04 12:15 | Tests written & model updated | Complete | Added Parent, ParentTask fields |
| 2026-03-04 12:20 | FormatTask updated | Complete | Shows parent in detail view |
| 2026-03-04 12:25 | getTask updated | Complete | Fetches parent task name |
| 2026-03-04 12:05 | All tests passing, build complete | Completed | make check & make build-all pass |
