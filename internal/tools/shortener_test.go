package tools

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGenerateToken(t *testing.T) {
	alphabetPattern := "[0-9a-zA-Z_]{10}"
	t.Log("Given the need to test regexp matching by generated tokens")
	for i := 0; i < 1000; i++ {
		token, err := GenerateToken()
		require.NoError(t, err, "Should be no error")
		require.Regexp(t, alphabetPattern, token, "Should match on regexp pattern")
	}
}
