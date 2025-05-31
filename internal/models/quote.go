package models

import "time"

type Quote struct {
    ID        int    `json:"id,omitempty"`
    Author    string    `json:"author"`
    Text      string    `json:"text"`
    CreatedAt time.Time `json:"created_at,omitempty"`
}