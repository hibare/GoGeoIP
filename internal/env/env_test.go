package env

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEnv(t *testing.T) {
	// Set Sample environment variables
	os.Setenv("STRING_ENV", "test_string")
	os.Setenv("BOOL_ENV", "true")
	os.Setenv("INT_ENV", "100")
	os.Setenv("DURATION_ENV", "24h")
	os.Setenv("SLICE_ENV", "1,2,3")

	Load()

	assert.Equal(t, "test_string", MustString("STRING_ENV", ""))
	assert.True(t, true, MustBool("BOOL_ENV", false))
	assert.Equal(t, 100, MustInt("INT_ENV", 0))
	assert.Equal(t, 24*time.Hour, MustDuration("DURATION_ENV", time.Duration(0)))
	assert.Equal(t, []string{"1", "2", "3"}, MustStringSlice("SLICE_ENV", []string{}))

	assert.Equal(t, "default_string", MustString("DEFAULT_STRING_ENV", "default_string"))
	assert.Equal(t, false, MustBool("DEFAULT_BOOL_ENV", false))
	assert.Equal(t, 1, MustInt("DEFAULT_INT_ENV", 1))
	assert.Equal(t, time.Duration(1), MustDuration("DEFAULT_DURATION_ENV", time.Duration(1)))
	assert.Equal(t, []string{"1"}, MustStringSlice("DEFAULT_SLICE_ENV", []string{"1"}))

	// Unset Sample environment variables
	os.Unsetenv("STRING_ENV")
	os.Unsetenv("BOOL_ENV")
	os.Unsetenv("INT_ENV")
	os.Unsetenv("DURATION_ENV")
	os.Unsetenv("SLICE_ENV")
}
