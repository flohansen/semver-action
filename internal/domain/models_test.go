package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion_IncByCommit(t *testing.T) {
	tests := []struct {
		name     string
		version  Version
		commit   Commit
		expected Version
	}{
		{
			name:     "breaking change increments major and resets minor and patch",
			version:  Version{Major: 1, Minor: 2, Patch: 3},
			commit:   Commit{IsBreaking: true},
			expected: Version{Major: 2, Minor: 0, Patch: 0},
		},
		{
			name:     "feat commit increments minor and resets patch",
			version:  Version{Major: 1, Minor: 2, Patch: 3},
			commit:   Commit{Type: "feat"},
			expected: Version{Major: 1, Minor: 3, Patch: 0},
		},
		{
			name:     "fix commit increments patch",
			version:  Version{Major: 1, Minor: 2, Patch: 3},
			commit:   Commit{Type: "fix"},
			expected: Version{Major: 1, Minor: 2, Patch: 4},
		},
		{
			name:     "chore commit increments patch",
			version:  Version{Major: 1, Minor: 2, Patch: 3},
			commit:   Commit{Type: "chore"},
			expected: Version{Major: 1, Minor: 2, Patch: 4},
		},
		{
			name:     "breaking change takes precedence over feat type",
			version:  Version{Major: 1, Minor: 2, Patch: 3},
			commit:   Commit{Type: "feat", IsBreaking: true},
			expected: Version{Major: 2, Minor: 0, Patch: 0},
		},
		{
			name:     "zero version incremented by breaking change",
			version:  Version{Major: 0, Minor: 0, Patch: 0},
			commit:   Commit{IsBreaking: true},
			expected: Version{Major: 1, Minor: 0, Patch: 0},
		},
		{
			name:     "zero version incremented by feat",
			version:  Version{Major: 0, Minor: 0, Patch: 0},
			commit:   Commit{Type: "feat"},
			expected: Version{Major: 0, Minor: 1, Patch: 0},
		},
		{
			name:     "zero version incremented by patch-level commit",
			version:  Version{Major: 0, Minor: 0, Patch: 0},
			commit:   Commit{Type: "fix"},
			expected: Version{Major: 0, Minor: 0, Patch: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			// act
			got := tt.version.IncByCommit(tt.commit)

			// assert
			assert.Equal(t, tt.expected, got)
			assert.NotEqual(t, got, tt.version)
		})
	}
}
