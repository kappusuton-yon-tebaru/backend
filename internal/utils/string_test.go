package utils_test

import (
	"testing"

	"github.com/kappusuton-yon-tebaru/backend/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestToKebabCase(t *testing.T) {
	tcs := []struct {
		name     string
		input    string
		expected string
	}{
		{
			"whitespace",
			"service name",
			"service-name",
		},
		{
			"snake case",
			"snake_case",
			"snake-case",
		},
		{
			"pascal case",
			"PascalCase",
			"pascalcase",
		},
		{
			"mixed",
			"PascalCase snake_case",
			"pascalcase-snake-case",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			kebab := utils.ToKebabCase(tc.input)
			assert.Equal(t, kebab, tc.expected)
		})
	}
}
