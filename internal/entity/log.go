package entity

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidLogLevel  = errors.New("invalid log level")
	ErrInvalidLogStatus = errors.New("invalid log status")
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

// Log represents a log
type Log struct {
	ID        string
	SystemID  string
	Level     string
	Status    string
	Message   string
	TimeStamp *time.Time
}

// NewLog creates a new Log
func NewLog(systemID, level, status, message string, timeStamp *time.Time) (*Log, error) {
	log := Log{
		ID:        uuid.New().String(),
		SystemID:  systemID,
		Level:     level,
		Status:    status,
		Message:   message,
		TimeStamp: timeStamp,
	}
	err := log.Validate()
	if err != nil {
		return nil, err
	}
	return &log, nil
}

// Validate checks if the log is valid. If the log is not valid, it returns an error specifying the invalidation.
func (l Log) Validate() error {
	_, err := uuid.Parse(l.ID)
	if err != nil {
		return err
	}
	_, err = uuid.Parse(l.SystemID)
	if err != nil {
		return err
	}

	if !isValidLogLevel(l.Level) {
		return ErrInvalidLogLevel
	}

	if !isValidLogStatus(l.Status) {
		return ErrInvalidLogStatus
	}

	if l.Message == "" {
		return ErrMessageRequired
	}

	if l.TimeStamp == nil {
		return ErrInvalidTimeStamp
	}

	return nil
}

func isValidLogLevel(level string) bool {
	_, ok := validLogLevels[strings.ToUpper(level)]
	return ok
}

func isValidLogStatus(status string) bool {
	_, ok := validLogStatus[strings.ToUpper(status)]
	return ok
}
