package security_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"alwis.dev/selectify/internal/security"
)

func TestNewToken_UniqueAndHashed(t *testing.T) {
	raw1, hash1, err := security.NewToken()
	require.NoError(t, err)
	raw2, hash2, err := security.NewToken()
	require.NoError(t, err)

	require.NotEmpty(t, raw1)
	require.NotEqual(t, raw1, raw2)
	require.NotEqual(t, hash1, hash2)
	require.Equal(t, security.HashToken(raw1), hash1)
	require.NotEqual(t, raw1, hash1)
}
