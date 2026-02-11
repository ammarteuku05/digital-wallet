package di

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Name string `validate:"required"`
}

func TestCustomValidator_Validate(t *testing.T) {
	v := NewCustomValidator()

	t.Run("valid struct", func(t *testing.T) {
		s := TestStruct{Name: "John"}
		err := v.Validate(s)
		assert.NoError(t, err)
	})

	t.Run("invalid struct", func(t *testing.T) {
		s := TestStruct{Name: ""}
		err := v.Validate(s)
		assert.Error(t, err)
	})
}
