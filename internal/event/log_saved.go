package event

import "time"

type LogSaved struct {
	Name    string
	Payload interface{}
}

func NewLogSaved(name string) *LogSaved {
	return &LogSaved{
		Name: name,
	}
}

func (e *LogSaved) GetName() string {
	return e.Name
}

func (e *LogSaved) GetPayload() interface{} {
	return e.Payload
}

func (e *LogSaved) SetPayload(payload interface{}) {
	e.Payload = payload
}

func (e *LogSaved) GetDateTime() time.Time {
	return time.Now()
}
