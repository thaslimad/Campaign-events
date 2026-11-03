# AI Usage Disclosure

## Tools used
- ChatGPT
- Go standard-library documentation (as needed)

## Use
ChatGPT assisted with interpreting the assignment, outlining the API, explaining validation/deduplication/concurrency decisions, drafting documentation, and reviewing the debugging exercise.

## Suggestion changed or rejected
An in-memory map was used for the first working backend because it avoids extra dependencies and supports the assignment's permitted in-memory option. For durable operation, I would replace it with SQLite and a unique event ID constraint. The in-memory version is not presented as production persistence.

## What AI helped me understand
The distinction between provider event IDs and contacts, the need to handle duplicates at concurrent request boundaries, and why timezone-dependent aggregation and unsynchronized counters are correctness bugs.

## AI-generated code that required debugging
Generated code was formatted, compiled, and checked against the provided debugging expected output. The final submission should only be represented as tested for commands actually run by the submitter.
