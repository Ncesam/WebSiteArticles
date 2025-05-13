package routers

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	authPb "backend/generated/proto/auth"
	"backend/internal/helpers"
	"backend/pkg/config"
	"backend/pkg/errors"
	jwt "backend/pkg/security/JWT"
	"backend/pkg/types"
)

// Register godoc
// @Summary      Регистрация
// @Description  Создание нового пользователя
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        data body types.RegisterForm true "Данные регистрации"
// @Success      200 {object} authPb.AuthResponse
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /auth/register [post]
func Register(logger *zap.Logger, cfg *config.Config, clients *types.MapClients) gin.HandlerFunc {
	return func(c *gin.Context) {
		var form types.RegisterForm
		if err := c.ShouldBindJSON(&form); err != nil {
			c.AbortWithStatusJSON(400, gin.H{"message": "Invalid data"})
			return
		}

		registerRequest := &authPb.RegisterRequest{
			Email:    form.Email,
			Nickname: form.Nickname,
			Password: form.Password,
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		_, err := clients.Auth.Service.Register(ctx, registerRequest)
		if err != nil {
			grpcErr := status.Convert(err)
			if grpcErr.Code() == codes.InvalidArgument {
				if grpcErr.Message() == "nickname, email and password are required" {
					c.AbortWithStatusJSON(errors.ErrInvalidCredentials.Code, gin.H{"message": "Fill fields"})
				} else if grpcErr.Message() == "password must be at least 8 characters" {
					c.AbortWithStatusJSON(errors.ErrInvalidCredentials.Code, gin.H{"message": "Password not validate"})
				}
			} else if grpcErr.Code() == codes.AlreadyExists {
				if grpcErr.Message() == "user with this email already exists" {
					c.AbortWithStatusJSON(errors.ErrUserAlreadyExists.Code, gin.H{"message": "email already exists"})
				} else if grpcErr.Message() == "user with this nickname already exists" {
					c.AbortWithStatusJSON(errors.ErrUserAlreadyExists.Code, gin.H{"message": "nickname already exists"})
				}
			} else {
				helpers.HandleGrpcError(logger, c, err, "Failed to register user")
			}
			return
		}

		c.JSON(200, gin.H{"message": "Registration successful"})
	}
}

// Login godoc
// @Summary      Авторизация
// @Description  Вход по email и паролю
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        data body types.LoginForm true "Данные входа"
// @Success      200 {object} authPb.AuthResponse
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Router       /auth/login [post]
func Login(logger *zap.Logger, cfg *config.Config, clients *types.MapClients) gin.HandlerFunc {
	return func(c *gin.Context) {
		var form types.LoginForm
		if err := c.ShouldBindJSON(&form); err != nil {
			c.AbortWithStatusJSON(400, gin.H{"message": "Invalid credentials"})
			return
		}

		req := &authPb.LoginRequest{
			Nickname: form.Nickname,
			Password: form.Password,
		}
		resp, err := clients.Auth.Service.Login(context.Background(), req)
		if err != nil {
			if status.Code(err) == codes.Unauthenticated {
				c.AbortWithStatusJSON(errors.ErrPermissionDenied.Code, gin.H{"message": "Password is invalid"})
			} else {
				helpers.HandleGrpcError(logger, c, err, "User not found")
			}
			return
		}
		c.SetCookie("refresh_token", resp.RefreshToken, 24*30*3600, "/", "", false, true)
		c.SetCookie("access_token", resp.AccessToken, 3600, "/", "", false, true)
		c.JSON(200, resp)
	}
}

// Logout godoc
// @Summary      Выход
// @Description  Инвалидация токена
// @Tags         auth
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]string
// @Router       /auth/logout [post]
func Logout(logger *zap.Logger, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetCookie("access_token", "", -1, "/", "", false, true)
		c.SetCookie("refresh_token", "", -1, "/", "", false, true)
		c.JSON(200, gin.H{"message": "Successfully logged out"})
	}
}

// Refresh godoc
// @Summary      Обновление токена
// @Description  Получить новый access токен
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200 {object} authPb.AuthResponse
// @Failure      401 {object} map[string]string
// @Router       /auth/refresh [post]
func Refresh(logger *zap.Logger, cfg *config.Config, clients *types.MapClients, authController *jwt.AuthController) gin.HandlerFunc {
	return func(c *gin.Context) {
		refreshToken, ok := helpers.GetRefreshToken(logger, c)
		if !ok {
			return
		}
		refreshTokenRequest := &authPb.RefreshRequest{
			RefreshToken: refreshToken,
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		resp, err := clients.Auth.Service.Refresh(ctx, refreshTokenRequest)
		if err != nil {
			grpcErr := status.Convert(err)
			if grpcErr.Code() == codes.InvalidArgument {
				if grpcErr.Message() == "refresh token is required" {
					c.AbortWithStatusJSON(errors.ErrInternalServer.Code, gin.H{"message": "refresh token is required"})
				}
			} else if grpcErr.Code() == codes.Unauthenticated {
				if grpcErr.Message() == "invalid refresh token" {
					c.AbortWithStatusJSON(errors.ErrPermissionDenied.Code, gin.H{"message": "refresh token is invaild"})
				}
			} else {
				helpers.HandleGrpcError(logger, c, err, "User not found")
			}
			return
		}
		c.SetCookie("refresh_token", resp.RefreshToken, 24*30*3600, "/", "", false, true)
		c.SetCookie("access_token", resp.AccessToken, 3600, "/", "", false, true)
		c.JSON(200, resp)
	}
}

// Me godoc
// @Summary      Получить информацию о себе
// @Description  Возвращает текущего пользователя из токена
// @Tags         auth
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} authPb.AuthResponse
// @Failure      401 {object} map[string]string
// @Router       /auth/me [put]
func Me(logger *zap.Logger, cfg *config.Config, clients *types.MapClients, authController *jwt.AuthController) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := helpers.CheckUser(logger, c, authController)
		if !ok {
			c.AbortWithStatusJSON(errors.ErrUnauthorized.Code, errors.ErrUnauthorized.Message)
			return
		}
		email, ok := claims["email"].(string)
		if !ok {
			c.AbortWithStatusJSON(errors.ErrInvalidEmailFormat.Code, errors.ErrInvalidEmailFormat.Message)
			return
		}
		getUserRequest := &authPb.GetMeRequest{
			Email: email,
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		resp, err := clients.Auth.Service.Me(ctx, getUserRequest)
		if err != nil {
			helpers.HandleGrpcError(logger, c, err, "User not found")
			return
		}
		c.SetCookie("refresh_token", resp.RefreshToken, 24*30*3600, "/", "", false, true)
		c.SetCookie("access_token", resp.AccessToken, 3600, "/", "", false, true)
		c.JSON(200, resp)
	}
}
