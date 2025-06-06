package jwt

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"backend/pkg/config"
	"backend/pkg/logger"
	"backend/pkg/types"
)

func TestAccessToken(t *testing.T) {
	data := []types.UserInfo{
		{1, "strtr", "fsdfsd"},
	}
	cfg, err := config.SetupConfig()
	require.NoError(t, err)

	loggerInstanse, err := logger.SetupLogger(logger.Parameters{zapcore.DebugLevel})
	require.NoError(t, err)

	AuthControllerInstanse := New(loggerInstanse, cfg)

	for _, item := range data {
		t.Run("", func(t *testing.T) {
			accessToken, err := AuthControllerInstanse.CreateAccessToken(item)
			require.NoError(t, err)

			mapClaims, err := AuthControllerInstanse.Decrypt(accessToken)
			require.NoError(t, err)

			id, ok := mapClaims["id"].(string)
			require.Equal(t, false, ok)

			email, ok := mapClaims["email"].(string)
			require.Equal(t, false, ok)

			username, ok := mapClaims["username"].(string)
			require.Equal(t, false, ok)

			assert.Equal(t, item.Id, id)
			assert.Equal(t, item.Email, email)
			assert.Equal(t, item.Nickname, username)
		})
	}

}
