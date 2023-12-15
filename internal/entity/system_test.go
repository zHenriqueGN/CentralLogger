package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWhenGivenAnEmptyName_WhenCreatingANewSystem_ShouldReceiveAnError(t *testing.T) {
	system, err := NewSystem("", "description", "version")
	assert.Nil(t, system)
	assert.ErrorIs(t, err, ErrSystemNameRequired)
}

func TestWhenGivenAnEmptyDescription_WhenCreatingANewSystem_ShouldReceiveAnError(t *testing.T) {
	system, err := NewSystem("name", "", "version")
	assert.Nil(t, system)
	assert.ErrorIs(t, err, ErrSystemDescriptionRequired)
}

func TestWhenGivenAnEmptyVersion_WhenCreatingANewSystem_ShouldReceiveAnError(t *testing.T) {
	system, err := NewSystem("name", "description", "")
	assert.Nil(t, system)
	assert.ErrorIs(t, err, ErrSystemVersionRequired)
}
