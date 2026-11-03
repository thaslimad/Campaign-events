# Campaign Events — Engineering Notes

## Interpretation
Relay receives provider webhooks for `sent`, `delivered`, `opened`, and `clicked`. The service accepts batches and returns per-campaign event counts. Provider retries must not inflate counts, and arrival order must not determine validity.

## Assumptions
- `event_id` is globally unique across providers and campaigns.
- The first valid payload for an event ID wins; later occurrences are duplicates.
- Counts represent accepted event occurrences, not unique contacts (except where specifically named unique opens in the separate debugging exercise).
- Timestamp must parse as RFC3339. The service preserves it and does not require events to arrive chronologically.
- A batch is processed item-by-item: malformed events are rejected while valid items continue.
- In this starter implementation, storage is in-memory and is lost on restart.

## Ambiguities / PM questions
- Is event ID globally unique or scoped to provider?
- If the same ID arrives with different payloads, should the first or latest payload win, or should this be surfaced?
- Are dashboard counts event counts or distinct-contact counts?
- Should APIs support date windows, pagination, and campaign existence checks?
- What is the retention and replay policy? Is eventual consistency acceptable?
- Should an invalid item cause whole-batch failure or partial acceptance?

## Priorities
1. Correct request decoding and per-event validation.
2. Deduplicate atomically while handling concurrent requests.
3. Provide campaign stats.
4. Exercise duplicate, invalid, and concurrency behavior.
5. Complete the debugging exercise and document scale considerations.

## Scale memo: 100 million events/day
The in-memory map will exhaust RAM, lose data on restart, and require a shared store across replicas. Scanning all events for each stats request will become too slow. First move durable ingestion to a database with a unique constraint on event ID and indexes on campaign/type/time; then introduce partitioning and asynchronous aggregation if measurements show it is needed. A queue can absorb bursts and decouple webhook acknowledgement from aggregation, but introduces delivery retries, ordering questions, consumer lag, poison messages, and at-least-once processing. Duplicates can occur at provider retry, queue redelivery, consumer restart, and database retry boundaries. Keep counters trustworthy with idempotent event persistence, transactional counter updates or rebuildable aggregates, reconciliation jobs, and monitoring for ingest rate, duplicate/reject rates, queue lag, DB latency, error rates, and aggregate drift. Keep deployment orchestration and multi-region complexity out until justified by measured requirements.

The external API can remain compatible, though pagination/date filtering and an explicit partial-batch response may be needed. Storage changes from process memory to durable indexed storage and likely a queue plus aggregate store. I would not introduce a complex streaming platform or multi-region design before bottlenecks are measured.

## Angry marketer: opened decreases from 250,000 to 248,000
Plausible explanations include a changed metric definition (event count vs unique contacts), a corrected or removed duplicate, late-arriving events that change a time-window boundary, a bot/privacy filtering adjustment, invalid-event quarantine, backfill/reprocessing, a campaign/timezone/filter mismatch, cache refresh, or a dashboard query/version change. Delivered rising by 5,000 is not inherently suspicious: delayed delivery events can arrive after the 10:00 snapshot.

First compare the exact query, campaign, time range, timezone, filters, metric definition, and data freshness for both snapshots. Inspect ingestion/rejection/de-duplication logs and any reprocessing or metric-definition deployment. Recompute the metric from raw accepted events for the same snapshot parameters. It is a bug if identical definitions and source data produce inconsistent results; it may be expected if the metric is mutable, windowed, or corrected.

## What I completed / what I intentionally skipped
Completed: standard-library Go ingestion and campaign stats endpoints, per-item validation, in-memory event-ID deduplication protected by a mutex, and the four fixes in the debugging exercise. Intentionally skipped: durable database, optional paginated events endpoint, authentication, deployment, and production queueing.

## If I had another day
I would use SQLite with a UNIQUE event ID constraint, add integration/load tests, define explicit batch semantics and metric definitions with the PM, and build a reconciliation path for aggregates.
