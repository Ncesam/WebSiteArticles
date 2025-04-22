package routers

import (
	"backend/internal/helpers"
	userPb "backend/generated/proto/user"
	"backend/pkg/config"
	"backend/pkg/errors"

	jwt "backend/pkg/security/JWT"
	"backend/pkg/types"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AddUserHandler godoc
// @Summary      Добавить пользователя
// @Description  Админ может создать нового пользователя
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        data body types.AddFormUser true "Пользователь для создания"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /user/add [post]
func AddUserHandler(logger *zap.Logger, cfg *config.Config, clients *types.MapClients, authController *jwt.AuthController) gin.HandlerFunc {
	return func(c *gin.Context) {
		var schema types.AddFormUser
		if err := c.ShouldBindJSON(&schema); err != nil {
			logger.Error("Request data are invalid", zap.Error(err))
			c.AbortWithStatusJSON(400, gin.H{"message": "Invalid request data"})
			return
		}

		if !helpers.IsAdmin(logger, c, authController) {
			return
		}

		addUserRequest := &userPb.AddUserRequest{
			Email:    schema.Email,
			Nickname: schema.Nickname,
			Password: schema.Password,
			IsAdmin:  schema.IsAdmin,
		}

		_, err := clients.User.Service.AddUser(clients.User.Ctx, addUserRequest)
		if err != nil {
			helpers.HandleGrpcError(logger, c, err, "Failed to add user")
			return
		}

		c.JSON(200, gin.H{"message": "User successfully added"})
	}
}

// GetUserHandler godoc
// @Summary      Получить информацию о пользователе
// @Description  Админ может указать nickname в query, обычный пользователь получает себя по токену
// @Tags         user
// @Produce      json
// @Security     BearerAuth
// @Param        nickname query string false "Никнейм пользователя (только для админов)"
// @Success      200 {object} userPb.User
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /user/get [get]
func GetUserHandler(logger *zap.Logger, cfg *config.Config, clients *types.MapClients, authController *jwt.AuthController) gin.HandlerFunc {
	return func(c *gin.Context) {
		var nickname string

		// Если админ — можно передать никнейм в query
		if helpers.IsAdmin(logger, c, authController) {
			nickname = c.Query("nickname")
			if nickname == "" {
				c.AbortWithStatusJSON(errors.ErrBadRequest.Code, gin.H{"message": "nickname is required"})
				return
			}
		} else {
			// Обычный пользователь — берем из токена
			claims, ok := helpers.CheckUser(logger, c, authController)
			if !ok {
				return
			}
			nicknameRaw, exists := claims["nickname"]
			if !exists {
				c.AbortWithStatusJSON(errors.ErrInvalidToken.Code, gin.H{"message": "nickname not found in token"})
				return
			}
			nickname, ok = nicknameRaw.(string)
			if !ok {
				c.AbortWithStatusJSON(errors.ErrInvalidToken.Code, gin.H{"message": "invalid nickname format"})
				return
			}
		}

		getUserRequest := &userPb.GetUserRequest{Nickname: nickname}
		response, err := clients.User.Service.GetUser(clients.User.Ctx, getUserRequest)
		if err != nil {
			helpers.HandleGrpcError(logger, c, err, "user not found")
			return
		}

		c.JSON(200, response)
	}
}

// UpdateUserHandler godoc
// @Summary      Обновить пользователя
// @Description  Админ может обновить любого, обычный — только себя
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        data body types.UpdateFormUser true "Данные для обновления"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Failure      403 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /user/update [put]
func UpdateUserHandler(logger *zap.Logger, cfg *config.Config, clients *types.MapClients, authController *jwt.AuthController) gin.HandlerFunc {
	return func(c *gin.Context) {
		var schema types.UpdateFormUser

		if err := c.ShouldBindJSON(&schema); err != nil {
			c.AbortWithStatusJSON(errors.ErrValidationFailed.Code, gin.H{"message": "Invalid update payload"})
			return
		}
		rawId, err := helpers.StringToInt32(logger, schema.Id)

		if err != nil {
			c.AbortWithStatusJSON(errors.ErrInvalidUserID.Code, gin.H{"message": errors.ErrInvalidUserID.Message})
			return
		}
		if !helpers.IsAdmin(logger, c, authController) {
			claims, ok := helpers.CheckUser(logger, c, authController)
			if !ok {
				return
			}

			id, ok := helpers.ExtractUserIDFromClaims(logger, c, claims)
			if !ok {
				return
			}

			if id != rawId {
				c.AbortWithStatusJSON(errors.ErrPermissionDenied.Code, gin.H{"message": errors.ErrPermissionDenied.Message})
				return
			}
		}

		updateUserRequest := &userPb.UpdateUserRequest{
			Id:       rawId,
			Nickname: schema.Nickname,
			Email:    schema.Email,
			Password: schema.Password,
			IsAdmin:  schema.IsAdmin,
		}

		_, err = clients.User.Service.UpdateUser(clients.User.Ctx, updateUserRequest)
		if err != nil {
			helpers.HandleGrpcError(logger, c, err, "Failed to update user")
			return
		}

		c.JSON(200, gin.H{"message": "User updated successfully"})
	}
}

// DeleteUserHandler godoc
// @Summary      Удалить пользователя
// @Description  Админ может удалить любого, пользователь — только себя
// @Tags         user
// @Produce      json
// @Security     BearerAuth
// @Param        id query int true "ID пользователя для удаления"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Failure      403 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /user/delete [delete]
func DeleteUserHandler(logger *zap.Logger, cfg *config.Config, clients *types.MapClients, authController *jwt.AuthController) gin.HandlerFunc {
	return func(c *gin.Context) {
		var deleteId int32
		var ok bool

		if helpers.IsAdmin(logger, c, authController) {
			deleteId, ok = helpers.ParseIDParam(logger, c, "id")
			if !ok {
				return
			}
		} else {
			claims, ok := helpers.CheckUser(logger, c, authController)
			if !ok {
				return
			}

			userId, ok := helpers.ExtractUserIDFromClaims(logger, c, claims)
			if !ok {
				return
			}

			deleteId, ok = helpers.ParseIDParam(logger, c, "id")
			if !ok {
				return
			}

			if userId != deleteId {
				c.AbortWithStatusJSON(errors.ErrPermissionDenied.Code, gin.H{"message": errors.ErrPermissionDenied.Message})
				return
			}
		}

		req := &userPb.DeleteUserRequest{
			Id: deleteId,
		}

		_, err := clients.User.Service.DeleteUser(clients.User.Ctx, req)
		if err != nil {
			helpers.HandleGrpcError(logger, c, err, "Failed to delete user")
			return
		}

		c.JSON(200, gin.H{"message": "User deleted successfully"})
	}
}
