package routers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"

	queuePb "backend/generated/proto/queue"
	"backend/internal/helpers"
	"backend/pkg/config"
	"backend/pkg/errors"
	jwt "backend/pkg/security/JWT"
	"backend/pkg/types"
)

// StartConfig godoc
// @Summary     Запуск конфигурации генерации
// @Description Принимает JSON с ConfigId, Prompt и Data и инициирует процесс генерации текста
// @Tags        config
// @Accept      json
// @Produce     json
// @Param       Authorization header string true "Bearer токен"
// @Param       body body types.InputForm true "Данные для генерации"
// @Success     200 {object} map[string]string "Successfully"
// @Failure     400 {object} map[string]string "Bad Request"
// @Failure     401 {object} map[string]string "Invalid Credentials"
// @Router      /queue/ [post]
func StartConfig(logger *zap.Logger, cfg *config.Config, clients *types.MapClients, authController *jwt.AuthController) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Debug("Handling StartConfig request")

		file, _, err := c.Request.FormFile("data")
		if err != nil {
			logger.Error("Invalid file", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid file"})
			return
		}

		prompt := c.PostForm("Prompt")
		userIDStr := c.PostForm("UserId")
		configID := c.PostForm("ConfigId")

		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid UserId"})
			return
		}
		body := types.InputRequestForm{
			Prompt:   prompt,
			UserId:   int64(userID),
			ConfigId: configID,
		}
		logger.Debug("Parsed input body", zap.String("config_id", body.ConfigId))
		workBook, err := excelize.OpenReader(file)
		if err != nil {
			logger.Error("Failed to read file", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to read Excel file"})
			return
		}

		sheetName := workBook.GetSheetName(0)
		if sheetName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "No sheet found in the Excel file"})
			return
		}

		rows, err := workBook.GetRows(sheetName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to read rows from the sheet"})
			return
		}

		if len(rows) < 1 {
			c.JSON(http.StatusOK, gin.H{"data": []interface{}{}})
			return
		}

		headers := rows[0]
		var data []map[string]string
		for _, row := range rows[1:] {
			item := make(map[string]string)
			for i, cell := range row {
				if i < len(headers) {
					item[headers[i]] = cell
				}
			}
			data = append(data, item)
		}
		jsonBytes, err := json.Marshal(data)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to encode data to JSON"})
			return
		}

		claims, ok := helpers.CheckUser(logger, c, authController)
		if !ok {
			logger.Error("User authentication failed", zap.String("error", errors.ErrInvalidCredentials.Message))
			c.AbortWithStatusJSON(errors.ErrInvalidCredentials.Code, gin.H{"message": errors.ErrInvalidCredentials.Message})
			return
		}

		userId, ok := helpers.ExtractUserIDFromClaims(logger, c, claims)
		if !ok {
			logger.Error("Failed to extract user ID from claims", zap.String("error", errors.ErrInvalidCredentials.Message))
			c.AbortWithStatusJSON(errors.ErrInvalidCredentials.Code, gin.H{"message": errors.ErrInvalidCredentials.Message})
			return
		}
		logger.Debug("Authenticated user", zap.Int64("user_id", userId))

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		logger.Debug("Sending gRPC StartConfig request", zap.String("config_id", body.ConfigId), zap.Int64("user_id", userId))
		_, err = clients.Queue.Service.StartConfig(ctx, &queuePb.StartConfigRequest{
			ConfigId: body.ConfigId,
			UserId:   userId,
			Prompt:   body.Prompt,
			Data:     string(jsonBytes),
		})
		if err != nil {
			logger.Error("StartConfig gRPC call failed", zap.Error(err))
			helpers.HandleGrpcError(logger, c, err, "Fail to start config")
			return
		}

		logger.Debug("Config started successfully")
		c.JSON(http.StatusOK, gin.H{"message": "Successfully"})
	}
}
