package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringInSlice(t *testing.T) {
	t.Run("string found in slice", func(t *testing.T) {
		slice := []string{"apple", "banana", "cherry"}
		result := StringInSlice("banana", slice)
		assert.True(t, result)
	})

	t.Run("string not found in slice", func(t *testing.T) {
		slice := []string{"apple", "banana", "cherry"}
		result := StringInSlice("grape", slice)
		assert.False(t, result)
	})

	t.Run("empty slice returns false", func(t *testing.T) {
		slice := []string{}
		result := StringInSlice("apple", slice)
		assert.False(t, result)
	})

	t.Run("empty string in slice", func(t *testing.T) {
		slice := []string{"", "apple", "banana"}
		result := StringInSlice("", slice)
		assert.True(t, result)
	})

	t.Run("case sensitive comparison", func(t *testing.T) {
		slice := []string{"Apple", "Banana"}
		result := StringInSlice("apple", slice)
		assert.False(t, result)
	})
}
