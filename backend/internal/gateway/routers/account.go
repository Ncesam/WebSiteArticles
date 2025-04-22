package routers

import (
	"backend/internal/helpers"
	accountPb "backend/generated/proto/account"
	"backend/pkg/config"
	"backend/pkg/errors"
	jwt "backend/pkg/security/JWT"
	"backend/pkg/types"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AddAccountHandler godoc
// @Summary      Добавить аккаунт
// @Description  Админ может создать новый аккаунт
// @Tags         account
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        data body types.AddAccountForm true "Аккаунт для создания"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /account/ [post]
func AddAccountHandler(logger *zap.Logger, cfg *config.Config, clients *types.MapClients, authController *jwt.AuthController) gin.HandlerFunc {
	return func(c *gin.Context) {
		var form types.AddAccountForm
		if err := c.ShouldBindJSON(&form); err != nil {
			logger.Error("Invalid AddAccountForm", zap.Error(err))
			c.AbortWithStatusJSON(400, gin.H{"message": "Invalid request data"})
			return
		}
		claims, ok := helpers.CheckUser(logger, c, authController)
		if !ok {
			return
		}

		userId, ok := helpers.ExtractUserIDFromClaims(logger, c, claims)
		if !ok {
			return
		}
		addAccountRequest := &accountPb.AddAccountRequest{
			Nickname: form.Nickname,
			Email:    form.Email,
			UserId:   userId,
		}
		_, err := clients.Account.Service.AddAccount(clients.Account.Ctx, addAccountRequest)
		if err != nil {
			helpers.HandleGrpcError(logger, c, err, "Failed to add account")
			return
		}

		c.JSON(200, gin.H{"message": "Account successfully added"})
	}
}

// ListAccountsHandler godoc
// @Summary      Получить список аккаунтов
// @Description  Получить список всех аккаунтов
// @Tags         account
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} accountPb.ListAccountResponse
// @Failure      500 {object} map[string]string
// @Router       /account/ [get]
func ListAccountsHandler(logger *zap.Logger, cfg *config.Config, clients *types.MapClients, authController *jwt.AuthController) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := helpers.CheckUser(logger, c, authController)
		if !ok {
			return
		}

		userId, ok := helpers.ExtractUserIDFromClaims(logger, c, claims)
		if !ok {
			return
		}

		listAccountRequest := &accountPb.ListAccountRequest{
			UserId: userId,
		}
		resp, err := clients.Account.Service.ListAccount(clients.Account.Ctx, listAccountRequest)
		if err != nil {
			helpers.HandleGrpcError(logger, c, err, "Failed to list accounts")
			return
		}
		c.JSON(200, resp)
	}
}

// DeleteAccountHandler godoc
// @Summary      Удалить аккаунт
// @Description  Админ может удалить аккаунт
// @Tags         account
// @Produce      json
// @Security     BearerAuth
// @Param        id query string true "ID аккаунта"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /account/ [delete]
func DeleteAccountHandler(logger *zap.Logger, cfg *config.Config, clients *types.MapClients, authController *jwt.AuthController) gin.HandlerFunc {
	return func(c *gin.Context) {

		idRaw := c.Query("id")
		if idRaw == "" {
			c.AbortWithStatusJSON(400, gin.H{"message": "id is required"})
			return
		}
		id, err := helpers.StringToInt32(logger, idRaw)
		if err != nil {
			c.AbortWithStatusJSON(errors.ErrBadRequest.Code, gin.H{"message":"Id field invalid"})
			return
		}
		deleteAccountRequest := &accountPb.DeleteAccountRequest{
			Id: id,
		}
		_, err = clients.Account.Service.DeleteAccount(clients.Account.Ctx, deleteAccountRequest)
		if err != nil {
			helpers.HandleGrpcError(logger, c, err, "Failed to delete account")
			return
		}

		c.JSON(200, gin.H{"message": "Account deleted successfully"})
	}
}
