package utils_test

import (
	"testing"

	"github.com/kappusuton-yon-tebaru/backend/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestParseVariable(t *testing.T) {
	tcs := []struct {
		name     string
		input    string
		expected []string
		err      error
	}{
		{
			name:     "empty string",
			input:    "",
			expected: nil,
			err:      utils.ErrBadVariable,
		},
		{
			name:     "empty variable",
			input:    "${}",
			expected: nil,
			err:      utils.ErrBadVariable,
		},
		{
			name:     "incomplete syntax",
			input:    "${",
			expected: nil,
			err:      utils.ErrBadVariable,
		},
		{
			name:     "single params",
			input:    "${service}",
			expected: []string{"service"},
			err:      nil,
		},
		{
			name:     "multiple params",
			input:    "${service::service_name::host}",
			expected: []string{"service", "service_name", "host"},
			err:      nil,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			results, err := utils.ParseVariable(tc.input)

			assert.Equal(t, tc.expected, results)
			assert.ErrorIs(t, err, tc.err)
		})
	}
}

func TestReplaceVariable(t *testing.T) {
	tcs := []struct {
		name     string
		input    map[string]string
		expected map[string]string
		err      error
	}{
		{
			name: "empty",
			input: map[string]string{
				"SERVICE_URL": "http://${}:3000",
			},
			expected: nil,
			err:      utils.ErrBadVariable,
		},
		{
			name: "single variable",
			input: map[string]string{
				"SERVICE_URL": "http://${HOST}:3000",
			},
			expected: map[string]string{
				"SERVICE_URL": "http://REPLACED_VARIABLE:3000",
			},
			err: nil,
		},
		{
			name: "multiple variables",
			input: map[string]string{
				"SERVICE_URL": "http://${HOST}:${PORT}",
			},
			expected: map[string]string{
				"SERVICE_URL": "http://REPLACED_VARIABLE:REPLACED_VARIABLE",
			},
			err: nil,
		},
		{
			name: "multiple variables invalid syntax",
			input: map[string]string{
				"SERVICE_URL": "http://${HOST}:${PORT",
			},
			expected: nil,
			err:      utils.ErrBadVariable,
		},
		{
			name: "multiple env vars",
			input: map[string]string{
				"SERVICE_URL": "http://${HOST}:${PORT}",
				"PORT":        "${PORT}",
			},
			expected: map[string]string{
				"SERVICE_URL": "http://REPLACED_VARIABLE:REPLACED_VARIABLE",
				"PORT":        "REPLACED_VARIABLE",
			},
			err: nil,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			replacedEnv, err := utils.ReplaceVariable(tc.input, func(params []string) (string, error) {
				return "REPLACED_VARIABLE", nil
			})

			assert.Equal(t, tc.expected, replacedEnv)
			assert.ErrorIs(t, err, tc.err)
		})
	}
}
