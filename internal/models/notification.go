package models

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	Id        uuid.UUID `json:"id"`
	Time      time.Time `json:"time"`
	OldStatus string    `json:"old_status"`
	NewStatus string    `json:"new_status"`
}
