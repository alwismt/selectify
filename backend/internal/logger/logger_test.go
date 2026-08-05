package logger

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestError(t *testing.T) {
	testError := fmt.Errorf("This is a test error %d", 1)
	err := Errorf(context.Background(), testError, "error message %d", 2)
	require.Error(t, err)
	require.Equal(t, "logger.TestError: This is a test error 1, error message 2", err.Error())

	err = Error(context.Background(), testError, "error message")
	require.Error(t, err)
	require.Equal(t, "logger.TestError: This is a test error 1, error message", err.Error())
}
