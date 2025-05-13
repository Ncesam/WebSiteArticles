package helpers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	jwtLib "github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"backend/pkg/errors"
	jwt "backend/pkg/security/JWT"
)

func IsAdmin(logger *zap.Logger, c *gin.Context, authController *jwt.AuthController) bool {
	accessToken, err := c.Cookie("access_token")
	if err != nil {
		logger.Warn("JWT not found in cookie", zap.Error(err))
		c.AbortWithStatusJSON(403, gin.H{"message": "JWT not found"})
		return false
	}

	claims, err := authController.Decrypt(accessToken)
	if err != nil {
		logger.Warn("Invalid JWT token during admin check", zap.Error(err))
		c.AbortWithStatusJSON(403, gin.H{"message": "Invalid JWT token"})
		return false
	}

	roleClaim, ok := claims["role"]
	if !ok {
		logger.Warn("Role claim not found in token")
		c.AbortWithStatusJSON(403, gin.H{"message": "role not found in token"})
		return false
	}

	role, ok := roleClaim.(bool)
	if !ok {
		logger.Error("Invalid role format in JWT", zap.Any("roleClaim", roleClaim))
		c.AbortWithStatusJSON(500, gin.H{"message": "Invalid role format"})
		return false
	}

	if !role {
		logger.Warn("Permission denied: user is not admin")
		c.AbortWithStatusJSON(403, gin.H{"message": "Permission denied"})
		return false
	}
	return true
}

func CheckUser(logger *zap.Logger, c *gin.Context, authController *jwt.AuthController) (jwtLib.MapClaims, bool) {
	accessToken, err := c.Cookie("access_token")
	if err != nil {
		logger.Warn("JWT not found in cookie", zap.Error(err))
		c.AbortWithStatusJSON(403, gin.H{"message": "JWT not found"})
		return nil, false
	}

	claims, err := authController.Decrypt(accessToken)
	if err != nil {
		logger.Warn("Failed to decrypt JWT", zap.Error(err))
		c.AbortWithStatusJSON(403, gin.H{"message": "Invalid JWT token"})
		return nil, false
	}
	return claims, true
}

func ParseIDParam(logger *zap.Logger, c *gin.Context, paramName string) (int64, bool) {
	idStr := c.Query(paramName)
	if idStr == "" {
		logger.Warn("Missing ID param", zap.String("param", paramName))
		c.AbortWithStatusJSON(errors.ErrBadRequest.Code, gin.H{"message": paramName + " is required"})
		return 0, false
	}

	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Warn("Invalid ID format", zap.String("param", paramName), zap.String("value", idStr), zap.Error(err))
		c.AbortWithStatusJSON(errors.ErrBadRequest.Code, gin.H{"message": "invalid " + paramName + " format"})
		return 0, false
	}

	return int64(idInt), true
}

func ExtractUserIDFromClaims(logger *zap.Logger, c *gin.Context, claims map[string]interface{}) (int64, bool) {
	subRaw, exists := claims["sub"]
	if !exists {
		logger.Warn("sub claim not found in token")
		c.AbortWithStatusJSON(errors.ErrInvalidToken.Code, gin.H{"message": "user ID not found in token"})
		return 0, false
	}

	var subStr string
	switch v := subRaw.(type) {
	case float64:
		subStr = strconv.Itoa(int(v))
	case string:
		subStr = v
	default:
		logger.Error("Invalid sub type", zap.Any("sub", subRaw))
		c.AbortWithStatusJSON(errors.ErrInvalidToken.Code, gin.H{"message": "invalid sub format"})
		return 0, false
	}

	idInt, err := strconv.Atoi(subStr)
	if err != nil {
		logger.Error("Failed to parse sub as int", zap.String("subStr", subStr), zap.Error(err))
		c.AbortWithStatusJSON(errors.ErrInvalidToken.Code, gin.H{"message": "invalid sub value"})
		return 0, false
	}

	return int64(idInt), true
}

func HandleGrpcError(logger *zap.Logger, c *gin.Context, err error, fallbackMessage string) {
	logger.Error(fallbackMessage, zap.Error(err))
	switch status.Code(err) {
	case codes.NotFound:
		c.AbortWithStatusJSON(errors.ErrNotFound.Code, gin.H{"message": fallbackMessage})
	case codes.PermissionDenied:
		c.AbortWithStatusJSON(errors.ErrPermissionDenied.Code, gin.H{"message": "Permission denied"})
	case codes.Unauthenticated:
		c.AbortWithStatusJSON(errors.ErrNotAuthed.Code, gin.H{"message": "Unauthenticated"})
	default:
		c.AbortWithStatusJSON(errors.ErrInternalServer.Code, gin.H{"message": "Internal server error"})
	}
}

func StringToInt64(logger *zap.Logger, s string) (int64, error) {
	i64, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		logger.Error("Failed to convert string to int64", zap.String("input", s), zap.Error(err))
	}
	return i64, err
}

func Int64ToString(i int64) string {
	return strconv.Itoa(int(i))
}

func GetRefreshToken(logger *zap.Logger, c *gin.Context) (string, bool) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		logger.Warn("Refresh token not found in cookie", zap.Error(err))
		c.AbortWithStatusJSON(errors.ErrTokenNotFound.Code, gin.H{"message": errors.ErrTokenNotFound.Message})
		return "", false
	}

	logger.Info("Refresh token retrieved successfully")
	return refreshToken, true
}

