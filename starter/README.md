# Campaign Events API

Requires Go 1.22+. Standard library only.

## Run

```bash
go run .
```

## Endpoints

`POST /events` accepts a JSON array. Each item is validated independently. The response reports `accepted`, `duplicates`, `rejected`, and optional per-item `errors`. Duplicate event IDs are ignored; first accepted payload wins. Malformed top-level JSON returns HTTP 400. The body limit is 10 MiB.

`GET /campaigns/{campaignID}/stats` returns counts of accepted event occurrences for sent, delivered, opened, and clicked.

## Example

```bash
curl -s -X POST localhost:8080/events -H 'Content-Type: application/json' -d '[{"event_id":"evt_1","campaign_id":"cmp_summer_sale","contact_id":"ct_001","type":"delivered","timestamp":"2026-08-10T06:15:00Z"}]'
curl -s localhost:8080/campaigns/cmp_summer_sale/stats
```

## Seed

```bash
curl -s -X POST localhost:8080/events -H 'Content-Type: application/json' --data-binary @seed/events.json
```

Storage is in-memory and resets on process restart. This implementation is a scoped assignment solution, not durable production storage.
