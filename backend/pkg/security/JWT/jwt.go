package jwt

import (
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"

	"backend/pkg/config"
	"backend/pkg/types"
)

type AccessToken = string
type RefreshToken = string

type AuthController struct {
	logger *zap.Logger
	cfg    *config.Config
}

func New(logger *zap.Logger, cfg *config.Config) AuthController {
	logger.Info("AuthController initialized")
	return AuthController{
		logger: logger,
		cfg:    cfg,
	}
}

func (controller AuthController) CreateAccessToken(user types.UserInfo) (AccessToken, error) {
	controller.logger.Info("Creating access token", zap.String("email", user.Email), zap.Int64("user_id", user.Id))

	claims := jwt.MapClaims{
		"nickname": user.Nickname,
		"email":    user.Email,
		"sub":      user.Id,
		"exp":      time.Now().Add(time.Duration(controller.cfg.AUTH.ACCESS.DURATION) * time.Second).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	keyHex := controller.cfg.AUTH.ACCESS.KEY
	keyBytes, err := hex.DecodeString(keyHex)
	if err != nil {
		controller.logger.Error("Failed to decode access token secret key", zap.Error(err))
		return "", err
	}

	signed, err := token.SignedString(keyBytes)
	if err != nil {
		controller.logger.Error("Failed to sign access token", zap.Error(err))
		return "", err
	}

	controller.logger.Debug("Access token created successfully")
	return AccessToken(signed), nil
}

func (controller AuthController) CreateRefreshToken(user types.UserInfo) (RefreshToken, error) {
	controller.logger.Info("Creating refresh token", zap.String("email", user.Email), zap.Int64("user_id", user.Id))

	claims := jwt.MapClaims{
		"nickname": user.Nickname,
		"email":    user.Email,
		"sub":      user.Id,
		"exp":      time.Now().Add(time.Duration(controller.cfg.AUTH.REFRESH.DURATION) * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	keyHex := controller.cfg.AUTH.REFRESH.KEY
	keyBytes, err := hex.DecodeString(keyHex)
	if err != nil {
		controller.logger.Error("Failed to decode refresh token secret key", zap.Error(err))
		return "", err
	}

	signed, err := token.SignedString(keyBytes)
	if err != nil {
		controller.logger.Error("Failed to sign refresh token", zap.Error(err))
		return "", err
	}

	controller.logger.Debug("Refresh token created successfully")
	return RefreshToken(signed), nil
}

func (controller AuthController) Decrypt(accessToken AccessToken) (jwt.MapClaims, error) {
	controller.logger.Debug("Decrypting access token")

	keyHex := controller.cfg.AUTH.ACCESS.KEY
	keyBytes, err := hex.DecodeString(keyHex)
	if err != nil {
		controller.logger.Error("Failed to decode access token key", zap.Error(err))
		return nil, err
	}

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			controller.logger.Error("Unexpected signing method", zap.String("alg", token.Method.Alg()))
			return nil, errors.New("unexpected signing method")
		}
		return keyBytes, nil
	})
	if err != nil {
		controller.logger.Error("Failed to parse access token", zap.Error(err))
		return nil, err
	}
	if !token.Valid {
		controller.logger.Error("Access token is invalid")
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		controller.logger.Error("Invalid claims type in access token")
		return nil, errors.New("invalid claims type")
	}

	exp, ok := claims["exp"].(float64)
	if !ok || time.Now().Unix() > int64(exp) {
		controller.logger.Error("Access token has expired", zap.Any("exp", claims["exp"]))
		return nil, errors.New("token expired")
	}

	controller.logger.Debug("Access token decrypted successfully")
	return claims, nil
}

func (controller AuthController) DecryptRefresh(refreshToken RefreshToken) (jwt.MapClaims, error) {
	controller.logger.Debug("Decrypting refresh token")

	keyHex := controller.cfg.AUTH.REFRESH.KEY
	keyBytes, err := hex.DecodeString(keyHex)
	if err != nil {
		controller.logger.Error("Failed to decode refresh token key", zap.Error(err))
		return nil, err
	}

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			controller.logger.Error("Unexpected signing method", zap.String("alg", token.Method.Alg()))
			return nil, errors.New("unexpected signing method")
		}
		return keyBytes, nil
	})
	if err != nil {
		controller.logger.Error("Failed to parse refresh token", zap.Error(err))
		return nil, err
	}
	if !token.Valid {
		controller.logger.Error("Refresh token is invalid")
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		controller.logger.Error("Invalid claims type in refresh token")
		return nil, errors.New("invalid claims type")
	}

	exp, ok := claims["exp"].(float64)
	if !ok || time.Now().Unix() > int64(exp) {
		controller.logger.Error("Refresh token has expired", zap.Any("exp", claims["exp"]))
		return nil, errors.New("token expired")
	}

	controller.logger.Debug("Refresh token decrypted successfully")
	return claims, nil
}
