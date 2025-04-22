package routers

import (
	queuePb "backend/generated/proto/queue"
	"backend/internal/helpers"
	"backend/pkg/config"
	"backend/pkg/types"
	"net/http"

	jwt "backend/pkg/security/JWT"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// AddQueueHandler godoc
// @Summary      Добавить элемент в очередь
// @Description  Аутентифицированный пользователь добавляет элемент в свою очередь
// @Tags         queue
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        data body int32 true "Данные для добавления в очередь"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /queue/ [post]
func AddQueueHandler(logger *zap.Logger, cfg *config.Config, clients *types.MapClients, authController *jwt.AuthController) gin.HandlerFunc {
	return func(c *gin.Context) {
		var priority int32
		if err := c.ShouldBindJSON(&priority); err != nil {
			logger.Error("Invalid AddQueueForm", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid request data"})
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

		req := &queuePb.AddQueueRequest{
			UserId:   userId,
			Priority: priority,
		}

		_, err := clients.Queue.Service.AddQueue(clients.Queue.Ctx, req)
		if err != nil {
			helpers.HandleGrpcError(logger, c, err, "Failed to add queue element")
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Queue element added"})
	}
}

// ListQueueHandler godoc
// @Summary      Получить элементы очереди
// @Description  Аутентифицированный пользователь получает все элементы своей очереди
// @Tags         queue
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} queuePb.ListQueueResponse
// @Failure      500 {object} map[string]string
// @Router       /queue/ [get]
func ListQueueHandler(logger *zap.Logger, cfg *config.Config, clients *types.MapClients, authController *jwt.AuthController) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := helpers.CheckUser(logger, c, authController)
		if !ok {
			return
		}
		userId, ok := helpers.ExtractUserIDFromClaims(logger, c, claims)
		if !ok {
			return
		}

		req := &queuePb.ListQueueRequest{UserId: userId}
		resp, err := clients.Queue.Service.ListQueue(clients.Queue.Ctx, req)
		if err != nil {
			helpers.HandleGrpcError(logger, c, err, "Failed to list queue elements")
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}

// ViewQueueHandler godoc
// @Summary      Просмотреть элемент очереди
// @Description  Получить информацию об элементе очереди по ID
// @Tags         queue
// @Produce      json
// @Security     BearerAuth
// @Param        id query int true "ID элемента очереди"
// @Success      200 {object} queuePb.ViewQueueResponse
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /queue/view [get]
func ViewQueueHandler(logger *zap.Logger, cfg *config.Config, clients *types.MapClients, authController *jwt.AuthController) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := helpers.ParseIDParam(logger, c, "id")
		if !ok {
			return
		}

		req := &queuePb.ViewQueueRequest{Id: id}
		resp, err := clients.Queue.Service.ViewQueue(clients.Queue.Ctx, req)
		if err != nil {
			helpers.HandleGrpcError(logger, c, err, "Failed to view queue element")
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}

// DelQueueHandler godoc
// @Summary      Удалить элемент очереди
// @Description  Удалить элемент очереди по ID
// @Tags         queue
// @Produce      json
// @Security     BearerAuth
// @Param        id query int true "ID элемента очереди"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /queue/ [delete]
func DeleteQueueHandler(logger *zap.Logger, cfg *config.Config, clients *types.MapClients, authController *jwt.AuthController) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := helpers.ParseIDParam(logger, c, "id")
		if !ok {
			return
		}

		_, err := clients.Queue.Service.DelQueue(clients.Queue.Ctx, &queuePb.DelQueueRequest{Id: id})
		if err != nil {
			helpers.HandleGrpcError(logger, c, err, "Failed to delete queue element")
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Queue element deleted"})
	}
}
