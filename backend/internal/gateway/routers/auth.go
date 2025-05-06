package routers

import (
	authPb "backend/generated/proto/auth"
	"backend/internal/helpers"
	"backend/pkg/config"
	jwt "backend/pkg/security/JWT"
	"backend/pkg/types"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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

		_, err := clients.Auth.Service.Register(clients.Auth.Ctx, registerRequest)
		if err != nil {
			helpers.HandleGrpcError(logger, c, err, "Failed to register user")
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

		resp, err := clients.Auth.Service.Login(clients.Auth.Ctx, req)
		if err != nil {
			helpers.HandleGrpcError(logger, c, err, "Login failed")
			return
		}

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
		resp, err := clients.Auth.Service.Refresh(clients.Auth.Ctx, refreshTokenRequest)
		if err != nil {
			helpers.HandleGrpcError(logger, c, err, "Token refresh failed")
			return
		}

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
// @Router       /auth/me [post]
func Me(logger *zap.Logger, cfg *config.Config, clients *types.MapClients, authController *jwt.AuthController) gin.HandlerFunc {
	return func(c *gin.Context) {
		refreshToken, ok := helpers.GetRefreshToken(logger, c)
		if !ok {
			return
		}
		getUserRequest := &authPb.GetMeRequest{
			RefreshToken: refreshToken,
		}

		resp, err := clients.Auth.Service.Me(clients.Auth.Ctx, getUserRequest)
		if err != nil {
			helpers.HandleGrpcError(logger, c, err, "Failed to get user data")
			return
		}
		c.JSON(200, resp)
	}
}
