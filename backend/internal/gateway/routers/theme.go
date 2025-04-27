package routers

import (
	themePb "backend/generated/proto/theme"
	"backend/internal/helpers"
	"backend/pkg/config"
	"backend/pkg/errors"
	jwt "backend/pkg/security/JWT"
	"backend/pkg/types"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AddThemeHandler(logger *zap.Logger, cfg *config.Config, clients *types.MapClients) gin.HandlerFunc {
	return func(c *gin.Context) {
		var form types.AddThemeForm
		if err := c.ShouldBindJSON(&form); err != nil {
			c.AbortWithStatusJSON(errors.ErrBadRequest.Code, gin.H{"message": "Invalid request"})
			return
		}
		addThemeRequest := &themePb.AddThemeRequest{
			Category: form.Category,
			Name: form.Name,
			Description: form.Description,
		}
		_, err := clients.Theme.Service.AddTheme(clients.Theme.Ctx, addThemeRequest)
		if err != nil {
			helpers.HandleGrpcError(logger, c, err, "failed to add theme")
			return
		}

		c.JSON(200, gin.H{"message": "Theme added"})
	}
}

func ListThemesHandler(logger *zap.Logger, cfg *config.Config, clients *types.MapClients, authController *jwt.AuthController) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := helpers.CheckUser(logger, c, authController)
		if !ok {
			return
		}

		userId, ok := helpers.ExtractUserIDFromClaims(logger, c, claims)
		if !ok {
			return
		}
		listThemesRequest := &themePb.ListThemesRequest{
			UserId: userId,
		}
		resp, err := clients.Theme.Service.ListThemes(clients.Theme.Ctx, listThemesRequest)
		if err != nil {
			helpers.HandleGrpcError(logger, c, err, "failed to list themes")
			return
		}
		c.JSON(200, resp)
	}
}

func ViewThemeHandler(logger *zap.Logger, cfg *config.Config, clients *types.MapClients) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := helpers.ParseIDParam(logger, c, "id")
		if !ok {
			return
		}
		viewThemeRequest := &themePb.ViewThemeRequest{
			Id: id,
		}
		resp, err := clients.Theme.Service.ViewTheme(clients.Theme.Ctx, viewThemeRequest)
		if err != nil {
			helpers.HandleGrpcError(logger, c, err, "failed to view theme")
			return
		}
		c.JSON(200, resp)
	}
}

func DeleteThemeHandler(logger *zap.Logger, cfg *config.Config, clients *types.MapClients) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := helpers.ParseIDParam(logger, c, "id")
		if !ok {
			return
		}

		deleteThemeRequest := &themePb.DeleteThemeRequest{
			Id: id,
		}

		_, err := clients.Theme.Service.DeleteTheme(clients.Theme.Ctx, deleteThemeRequest)
		if err != nil {
			helpers.HandleGrpcError(logger, c, err, "failed to delete theme")
			return
		}
		c.JSON(200, gin.H{"message": "Theme deleted"})
	}
}
