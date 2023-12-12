package entity

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Log represents a log
type Log struct {
	ID        uuid.UUID  `json:"id"`
	SystemID  uuid.UUID  `json:"system_id"`
	Level     string     `json:"level"`
	Status    string     `json:"status"`
	Message   string     `json:"message"`
	TimeStamp *time.Time `json:"time_stamp"`
	UserID    uuid.UUID  `json:"user_id"`
}

// New creates a new Log
func New(systemID uuid.UUID, level string, status string, message string, timeStamp *time.Time, userID uuid.UUID) (*Log, error) {
	log := Log{
		ID:        uuid.New(),
		SystemID:  systemID,
		Level:     level,
		Status:    status,
		Message:   message,
		TimeStamp: timeStamp,
		UserID:    userID,
	}
	err := log.IsValid()
	if err != nil {
		return nil, err
	}
	return &log, nil
}

var (
	ErrInvalidLogLevel  = errors.New("invalid log level")
	ErrMessageRequired  = errors.New("message is required")
	ErrInvalidTimeStamp = errors.New("invalid time stamp")
)

var validLogLevels = map[string]bool{
	"TRACE":     true,
	"DEBUG":     true,
	"INFO":      true,
	"NOTICE":    true,
	"WARNING":   true,
	"ERROR":     true,
	"CRITICAL":  true,
	"ALERT":     true,
	"EMERGENCY": true,
}

var validLogStatus = map[string]bool{
	"SUCCESS": true,
	"FAILURE": true,
}

func isValidLogLevel(level string) bool {
	_, ok := validLogLevels[strings.ToUpper(level)]
	return ok
}

func isValidLogStatus(status string) bool {
	_, ok := validLogStatus[strings.ToUpper(status)]
	return ok
}

// IsValid checks if the log is valid. If the log is not valid, it returns an error specifying the invalidation.
func (l Log) IsValid() error {
	_, err := uuid.Parse(l.ID.String())
	if err != nil {
		return err
	}
	_, err = uuid.Parse(l.SystemID.String())
	if err != nil {
		return err
	}

	if !isValidLogLevel(l.Level) {
		return ErrInvalidLogLevel
	}

	if !isValidLogStatus(l.Status) {
		return ErrInvalidLogLevel
	}

	if l.Message == "" {
		return ErrMessageRequired
	}

	if l.TimeStamp == nil {
		return ErrInvalidTimeStamp
	}

	_, err = uuid.Parse(l.UserID.String())
	if err != nil {
		return err
	}
	return nil
}
