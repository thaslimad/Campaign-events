package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	mu     sync.RWMutex
	events = make(map[string]Event)
)

var validTypes = map[string]bool{"sent": true, "delivered": true, "opened": true, "clicked": true}

type BatchResponse struct {
	Accepted   int      `json:"accepted"`
	Duplicates int      `json:"duplicates"`
	Rejected   int      `json:"rejected"`
	Errors     []string `json:"errors,omitempty"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /events", handlePostEvents)
	mux.HandleFunc("GET /campaigns/{campaignID}/stats", handleGetStats)
	log.Printf("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func handlePostEvents(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSON(w, 400, map[string]string{"error": "request body too large or unreadable"})
		return
	}
	var raw []json.RawMessage
	if err = json.Unmarshal(body, &raw); err != nil || raw == nil {
		writeJSON(w, 400, map[string]string{"error": "request must be a JSON array"})
		return
	}
	result := BatchResponse{}
	for i, item := range raw {
		var ev Event
		if err := json.Unmarshal(item, &ev); err != nil {
			result.Rejected++
			result.Errors = append(result.Errors, fmt.Sprintf("event %d: invalid JSON event", i))
			continue
		}
		if err := validateEvent(ev); err != nil {
			result.Rejected++
			result.Errors = append(result.Errors, fmt.Sprintf("event %d: %s", i, err))
			continue
		}
		mu.Lock()
		if _, ok := events[ev.EventID]; ok {
			mu.Unlock()
			result.Duplicates++
			continue
		}
		events[ev.EventID] = ev
		mu.Unlock()
		result.Accepted++
	}
	writeJSON(w, http.StatusOK, result)
}

func validateEvent(ev Event) error {
	if strings.TrimSpace(ev.EventID) == "" {
		return fmt.Errorf("event_id is required")
	}
	if strings.TrimSpace(ev.CampaignID) == "" {
		return fmt.Errorf("campaign_id is required")
	}
	if strings.TrimSpace(ev.ContactID) == "" {
		return fmt.Errorf("contact_id is required")
	}
	if !validTypes[ev.Type] {
		return fmt.Errorf("invalid event type")
	}
	if _, err := time.Parse(time.RFC3339, ev.Timestamp); err != nil {
		return fmt.Errorf("timestamp must be RFC3339")
	}
	if len(ev.Metadata) > 0 && !json.Valid(ev.Metadata) {
		return fmt.Errorf("metadata must be valid JSON")
	}
	return nil
}

func handleGetStats(w http.ResponseWriter, r *http.Request) {
	campaignID := r.PathValue("campaignID")
	s := Stats{CampaignID: campaignID}
	mu.RLock()
	for _, ev := range events {
		if ev.CampaignID != campaignID {
			continue
		}
		switch ev.Type {
		case "sent":
			s.Sent++
		case "delivered":
			s.Delivered++
		case "opened":
			s.Opened++
		case "clicked":
			s.Clicked++
		}
	}
	mu.RUnlock()
	writeJSON(w, 200, s)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("writeJSON: %v", err)
	}
}
