package models

import (
	"encoding/json"
	"time"
)

type Job struct {
	ID        int64           `json:"id" db:"id"`
	Payload   json.RawMessage `json:"payload" db:"payload"`
	Status    string          `json:"status" db:"status"`
	Result    json.RawMessage `json:"result" db:"result"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" db:"updated_at"`
}
