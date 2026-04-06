package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeCommitFromString(t *testing.T) {
	for _, tt := range []struct {
		name     string
		message  string
		expected Commit
	}{
		{
			name:    "",
			message: `feat: add new feature`,
			expected: Commit{
				Raw:        "feat: add new feature",
				Type:       "feat",
				Scope:      "",
				Message:    "add new feature",
				IsBreaking: false,
			},
		},
		{
			name:    "",
			message: `feat(test): add new feature`,
			expected: Commit{
				Raw:        "feat(test): add new feature",
				Type:       "feat",
				Scope:      "test",
				Message:    "add new feature",
				IsBreaking: false,
			},
		},
		{
			name:    "",
			message: `feat!: add new feature`,
			expected: Commit{
				Raw:        "feat!: add new feature",
				Type:       "feat",
				Scope:      "",
				Message:    "add new feature",
				IsBreaking: true,
			},
		},
		{
			name:    "",
			message: `feat(test)!: add new feature`,
			expected: Commit{
				Raw:        "feat(test)!: add new feature",
				Type:       "feat",
				Scope:      "test",
				Message:    "add new feature",
				IsBreaking: true,
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			// act
			commit, err := DecodeCommitFromString(tt.message)

			// assert
			require.NoError(t, err)
			assert.Equal(t, tt.expected, commit)
		})
	}
}

func TestDecodeVersionFromString(t *testing.T) {
	for _, tt := range []struct {
		name     string
		version  string
		expected Version
	}{
		{
			name:    "",
			version: "1.2.3",
			expected: Version{
				Major: 1,
				Minor: 2,
				Patch: 3,
			},
		},
		{
			name:    "",
			version: "v1.2.3",
			expected: Version{
				Major: 1,
				Minor: 2,
				Patch: 3,
			},
		},
		{
			name:    "",
			version: "v10.20.30",
			expected: Version{
				Major: 10,
				Minor: 20,
				Patch: 30,
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			// act
			v, err := DecodeVersionFromString(tt.version)

			// assert
			require.NoError(t, err)
			assert.Equal(t, tt.expected, v)
		})
	}
}
