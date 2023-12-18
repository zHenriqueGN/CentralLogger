package event

import "time"

type SystemCreated struct {
	Name    string
	Payload interface{}
}

func NewSystemCreated(name string) *SystemCreated {
	return &SystemCreated{
		Name: name,
	}
}

func (s *SystemCreated) GetName() string {
	return s.Name
}

func (s *SystemCreated) GetPayload() interface{} {
	return s.Payload
}

func (s *SystemCreated) SetPayLoad(payload interface{}) {
	s.Payload = payload
}

func (s *SystemCreated) GetDateTime() time.Time {
	return time.Now()
}
