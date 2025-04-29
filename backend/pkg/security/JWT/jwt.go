package jwt

import (
	"backend/pkg/config"
	"backend/pkg/types"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
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
		"exp":      time.Now().Add(time.Duration(controller.cfg.AUTH.ACCESS.DURATION)).Unix(),
		"iat":      time.Now().Unix(),
	}
	controller.logger.Debug("Create New JWT")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString(controller.cfg.AUTH.ACCESS.KEY)
	if err != nil {
		controller.logger.Error("Error in crypting token")
		return "", err
	}

	return AccessToken(signed), nil
}

func (controller AuthController) Decrypt(accessToken AccessToken) (jwt.MapClaims, error) {
	token, err := jwt.Parse(string(accessToken), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			controller.logger.Error("Unexpected signing method")
			return nil, errors.New("unexpected signing method")
		}
		return []byte(controller.cfg.AUTH.ACCESS.KEY), nil
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

	expTime := time.Unix(int64(exp), 0)
	if !expTime.Before(time.Now()) {
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
		"exp":      time.Now().Add(time.Duration(controller.cfg.AUTH.REFRESH.DURATION)).Unix(),
		"iat":      time.Now().Unix(),
	}
	controller.logger.Debug("Create New JWT")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString(controller.cfg.AUTH.REFRESH.KEY)
	if err != nil {
		controller.logger.Error("Error in crypting token")
		return "", err
	}

	return RefreshToken(signed), nil
}
