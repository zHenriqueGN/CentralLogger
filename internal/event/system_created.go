package event

import "time"

type SystemCreated struct {
	Name    string
	Payload interface{}
}

func NewSystemCreated() *SystemCreated {
	return &SystemCreated{
		Name: "SystemCreated",
	}
}

func (s *SystemCreated) GetName() string {
	return s.Name
}

func (s *SystemCreated) GetPayload() interface{} {
	return s.Payload
}

func (s *SystemCreated) SetPayload(payload interface{}) {
	s.Payload = payload
}

func (s *SystemCreated) GetDateTime() time.Time {
	return time.Now()
}
