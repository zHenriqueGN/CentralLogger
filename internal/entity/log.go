package entity

import (
	"time"

	"github.com/google/uuid"
)

type Log struct {
	ID        uuid.UUID  `json:"id"`
	SystemID  uuid.UUID  `json:"system_id"`
	Severity  string     `json:"severity"`
	Status    string     `json:"status"`
	Message   string     `json:"message"`
	TimeStamp *time.Time `json:"time_stamp"`
	UserID    uuid.UUID  `json:"user_id"`
}
