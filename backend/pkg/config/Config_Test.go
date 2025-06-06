package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	cfg, err := SetupConfig()
	assert.NoError(t, err)

	assert.Equal(t, "localhost", cfg.APP.HOST)
	assert.Equal(t, 8000, cfg.APP.PORT)

	assert.Equal(t, 120, cfg.AUTH.ACCESS.DURATION)
	assert.NotEqual(t, "secretkey", cfg.AUTH.ACCESS.KEY)

	assert.Equal(t, 120, cfg.AUTH.REFRESH.DURATION)
	assert.NotEqual(t, "secretkey", cfg.AUTH.REFRESH.KEY)

	assert.NotEqual(t, "testuser", cfg.MONGO.USERNAME)
	assert.NotEqual(t, "testpassword", cfg.MONGO.PASSWORD)
	assert.NotEqual(t, "testuser", cfg.POSTGRES.USERNAME)
	assert.NotEqual(t, "testpassword", cfg.POSTGRES.PASSWORD)
}
