package entity

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGivenAnInvalidSystemID_WhenCreatingANewLog_ThenShouldReceiveAnError(t *testing.T) {
	currentTime := time.Now()
	log, err := NewLog("wrong_uuid", "WRONG_LEVEL", "SUCCESS", "MESSAGE", &currentTime)
	assert.Nil(t, log)
	assert.Error(t, err)
}

func TestGivenAnInvalidLogLevel_WhenCreatingANewLog_ThenShouldReceiveAnError(t *testing.T) {
	currentTime := time.Now()
	log, err := NewLog(uuid.New().String(), "WRONG_LEVEL", "SUCCESS", "MESSAGE", &currentTime)
	assert.Nil(t, log)
	assert.ErrorIs(t, err, ErrInvalidLogLevel)
}

func TestGivenAnInvalidLogStatus_WhenCreatingANewLog_ThenShouldReceiveAnError(t *testing.T) {
	currentTime := time.Now()
	log, err := NewLog(uuid.New().String(), "DEBUG", "WRONG_STATUS", "MESSAGE", &currentTime)
	assert.Nil(t, log)
	assert.ErrorIs(t, err, ErrInvalidLogStatus)
}

func TestGivenAnEmptyMessage_WhenCreatingANewLog_ThenShouldReceiveAnError(t *testing.T) {
	currentTime := time.Now()
	log, err := NewLog(uuid.New().String(), "DEBUG", "SUCCESS", "", &currentTime)
	assert.Nil(t, log)
	assert.ErrorIs(t, err, ErrMessageRequired)
}

func TestGivenAnNilTimeStamp_WhenCreatingANewLog_ThenShouldReceiveAnError(t *testing.T) {
	log, err := NewLog(uuid.New().String(), "DEBUG", "SUCCESS", "MESSAGE", nil)
	assert.Nil(t, log)
	assert.ErrorIs(t, err, ErrInvalidTimeStamp)
}

func TestGivenValidFields_WhenCreatingANewLog_ThenShouldReceiveTheLog(t *testing.T) {
	currentTime := time.Now()
	log, err := NewLog(uuid.New().String(), "DEBUG", "SUCCESS", "MESSAGE", &currentTime)
	assert.Nil(t, err)
	assert.NotNil(t, log)
	assert.Len(t, log.ID, 36)
	assert.Equal(t, log.Level, "DEBUG")
	assert.Equal(t, log.Status, "SUCCESS")
	assert.Equal(t, log.Message, "MESSAGE")
	assert.Equal(t, *log.TimeStamp, currentTime)
}
