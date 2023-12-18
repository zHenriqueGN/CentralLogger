package event

import "time"

type LogSaved struct {
	Name    string
	Payload string
}

func NewLogSaved(name string) *LogSaved {
	return &LogSaved{
		Name: name,
	}
}

func (e *LogSaved) GetName() string {
	return e.Name
}

func (e *LogSaved) GetPayload() string {
	return e.Payload
}

func (e *LogSaved) SetPayload(payload string) {
	e.Payload = payload
}

func (e *LogSaved) GetDateTime(name string) time.Time {
	return time.Now()
}
