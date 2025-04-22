package routers

import (
	"backend/internal/helpers"
	articlePb "backend/generated/proto/article"
	"backend/pkg/config"
	"backend/pkg/errors"
	jwt "backend/pkg/security/JWT"
	"backend/pkg/types"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AddArticleHandler godoc
// @Summary      Добавить статью
// @Description  Только админ может добавлять статьи
// @Tags         article
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        data body types.AddArticleForm true "Статья для добавления"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Failure      403 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /articles/ [post]
func AddArticleHandler(logger *zap.Logger, cfg *config.Config, clients *types.MapClients, authController *jwt.AuthController) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.AddArticleForm

		if err := c.ShouldBindJSON(&req); err != nil {
			logger.Error("Invalid article data", zap.Error(err))
			c.AbortWithStatusJSON(errors.ErrValidationFailed.Code, gin.H{"message": "Invalid request data"})
			return
		}

		_, ok := helpers.CheckUser(logger, c, authController)
		if !ok {
			return
		}

		article := &articlePb.AddArticleRequest{
			Title:   req.Title,
			Content: req.Content,
		}

		_, err := clients.Article.Service.AddArticle(clients.Article.Ctx, article)
		if err != nil {
			helpers.HandleGrpcError(logger, c, err, "Failed to add article")
			return
		}

		c.JSON(200, gin.H{"message": "Article successfully added"})
	}
}

// ListArticlesHandler godoc
// @Summary      Получить список статей
// @Description  Все пользователи могут получить список статей
// @Tags         article
// @Produce      json
// @Success      200 {array} articlePb.ListArticlesResponse
// @Failure      500 {object} map[string]string
// @Router       /articles/ [get]
func ListArticlesHandler(logger *zap.Logger, cfg *config.Config, clients *types.MapClients, authController *jwt.AuthController) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, ok := helpers.CheckUser(logger, c, authController)
		if !ok {
			return
		}
		resp, err := clients.Article.Service.ListArticles(clients.Article.Ctx, &articlePb.ListArticlesRequest{})
		if err != nil {
			helpers.HandleGrpcError(logger, c, err, "Failed to list articles")
			return
		}

		c.JSON(200, resp)
	}
}
