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
	return AuthController{
		logger: logger, cfg: cfg,
	}
}
func (controller AuthController) CreateAccessToken(user types.UserInfo) (AccessToken, error) {
	claims := jwt.MapClaims{
		"nickname": user.Nickname,
		"email":    user.Email,
		"sub":      user.Id,
		"exp":      time.Now().Add(time.Duration(controller.cfg.AUTH.ACCESS.DURATION) * time.Second).Unix(),
		"iat":      time.Now().Unix(),
	}
	controller.logger.Debug("Create New JWT")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	keyHex := controller.cfg.AUTH.ACCESS.KEY
	keyBytes, err := hex.DecodeString(keyHex)
	if err != nil {
		controller.logger.Error("Failed to decode secret key")
		return "", err
	}
	signed, err := token.SignedString(keyBytes)
	if err != nil {
		controller.logger.Error("Error in crypting token")
		return "", err
	}

	return AccessToken(signed), nil
}
func (controller AuthController) Decrypt(accessToken AccessToken) (jwt.MapClaims, error) {
	keyHex := controller.cfg.AUTH.ACCESS.KEY
	keyBytes, err := hex.DecodeString(keyHex)
	if err != nil {
		controller.logger.Error("Failed to decode secret key")
		return nil, err
	}

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("unexpected signing method")
		}
		return keyBytes, nil
	})


	if err != nil || !token.Valid {
		controller.logger.Error("Token is invalid")
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		controller.logger.Error("Invalid claims type")
		return nil, errors.New("invalid claims type")
	}
	exp, ok := claims["exp"].(float64)
	if !ok {
		controller.logger.Error("Token expired")
		return nil, errors.New("token expired")
	}

	if time.Now().Unix() > int64(exp){
		controller.logger.Error("Token expired")
		return nil, errors.New("token expired")
	}

	return claims, nil
}
func (controller AuthController) DecryptRefresh(refreshToken RefreshToken) (jwt.MapClaims, error) {
	keyHex := controller.cfg.AUTH.REFRESH.KEY
	keyBytes, err := hex.DecodeString(keyHex)
	if err != nil {
		controller.logger.Error("Failed to decode secret key")
		return nil, err
	}

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("unexpected signing method")
		}
		return keyBytes, nil
	})

	if err != nil || !token.Valid {
		controller.logger.Error("Token is invalid")
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		controller.logger.Error("Invalid claims type")
		return nil, errors.New("invalid claims type")
	}
	exp, ok := claims["exp"].(float64)
	if !ok {
		controller.logger.Error("Token expired")
		return nil, errors.New("token expired")
	}

	if time.Now().Unix() > int64(exp){
		controller.logger.Error("Token expired")
		return nil, errors.New("token expired")
	}

	return claims, nil
}
func (controller AuthController) CreateRefreshToken(user types.UserInfo) (RefreshToken, error) {
	claims := jwt.MapClaims{
		"nickname": user.Nickname,
		"email":    user.Email,
		"sub":      user.Id,
		"exp":      time.Now().Add(time.Duration(controller.cfg.AUTH.REFRESH.DURATION) * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}
	controller.logger.Debug("Create New JWT")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	keyHex := controller.cfg.AUTH.ACCESS.KEY
	keyBytes, err := hex.DecodeString(keyHex)
	if err != nil {
		controller.logger.Error("Failed to decode secret key")
		return "", err
	}
	signed, err := token.SignedString(keyBytes)
	if err != nil {
		controller.logger.Error("Error in crypting token")
		return "", err
	}

	return RefreshToken(signed), nil
}
