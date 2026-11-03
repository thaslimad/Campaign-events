# Campaign Events Engineering Assignment

Go HTTP service for provider webhook ingestion and campaign statistics, plus the corrected debugging exercise.

## Run the backend

```bash
cd starter
go run .
```

Endpoints:
- `POST /events` — JSON array; responds with accepted, duplicate, and rejected counts.
- `GET /campaigns/{campaignID}/stats` — counts accepted events by type.

Example:

```bash
curl -X POST localhost:8080/events -H 'Content-Type: application/json' -d '[{"event_id":"evt_1","campaign_id":"cmp_demo","contact_id":"ct_1","type":"sent","timestamp":"2026-08-10T06:15:00Z"}]'
curl localhost:8080/campaigns/cmp_demo/stats
```

Run Go checks:

```bash
go test ./...
go vet ./...
go run -race .
```

## Debugging exercise

```bash
cd debugging
go run . events.jsonl > actual.txt
diff actual.txt expected_output.txt
go run -race . events.jsonl
```

See `NOTES.md`, `BUGS.md`, and `AI_USAGE.md` for decisions and disclosures. The backend uses in-memory storage; data resets on restart.
