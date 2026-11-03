# Debugging findings

## Bug 1 — Deduplication only within each batch
- **What:** `seen` was allocated inside `processBatch`, so a repeated event ID in a later batch was counted again.
- **Why:** The input is processed in multiple batches, but the spec requires uniqueness across the whole file.
- **Fix:** Use package-level `seenEventIDs` for the entire run.
- **Verification:** Ran `go run . events.jsonl` and diffed stdout against `expected_output.txt`; diff was empty.

## Bug 2 — Unique opens keyed only by contact
- **What:** `openedBy` used only `ContactID`.
- **Why:** The same contact opening in different campaigns should count once for each campaign, but the global contact key suppressed the second campaign.
- **Fix:** Key by campaign ID plus contact ID.
- **Verification:** Exact output matched the expected output.

## Bug 3 — Daily delivered bucket used local date
- **What:** `Timestamp.Local()` selected the host's timezone.
- **Why:** The requirements specify UTC dates, and output could vary with machine timezone.
- **Fix:** Use `Timestamp.UTC()` before formatting the date.
- **Verification:** Exact output matched expected output; UTC conversion is independent of host timezone.

## Bug 4 — Concurrent counter increments raced
- **What:** Worker goroutines incremented shared plain integer fields without synchronization.
- **Why:** Multiple workers can update one campaign's stats concurrently.
- **Fix:** Add `applyMu` around the counter update in `apply`.
- **Verification:** `go run -race . events.jsonl` and repeated exact-output comparisons.
