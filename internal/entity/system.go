package entity

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrSystemNameRequired        = errors.New("system name is required")
	ErrSystemDescriptionRequired = errors.New("system description is required")
	ErrSystemVersionRequired     = errors.New("system version is required")
)

// System represents a system that generates logs
type System struct {
	ID          uuid.UUID
	Name        string
	Description string
	Version     string
}

// NewSystem creates a new System
func NewSystem(name, description, version string) (*System, error) {
	system := System{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		Version:     version,
	}
	err := system.Validate()
	if err != nil {
		return nil, err
	}
	return &system, nil
}

// Validate checks if the log is valid. If the log is not valid, it returns an error specifying the invalidation.
func (s *System) Validate() error {
	_, err := uuid.Parse(s.ID.String())
	if err != nil {
		return err
	}
	if s.Name == "" {
		return ErrSystemNameRequired
	}
	if s.Description == "" {
		return ErrSystemDescriptionRequired
	}
	if s.Version == "" {
		return ErrSystemVersionRequired
	}
	return nil
}
