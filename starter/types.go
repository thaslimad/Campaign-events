package main

import "encoding/json"

type Event struct {
	EventID    string          `json:"event_id"`
	CampaignID string          `json:"campaign_id"`
	ContactID  string          `json:"contact_id"`
	Type       string          `json:"type"`
	Timestamp  string          `json:"timestamp"`
	Metadata   json.RawMessage `json:"metadata,omitempty"`
}

type Stats struct {
	CampaignID string `json:"campaign_id"`
	Sent       int    `json:"sent"`
	Delivered  int    `json:"delivered"`
	Opened     int    `json:"opened"`
	Clicked    int    `json:"clicked"`
}
