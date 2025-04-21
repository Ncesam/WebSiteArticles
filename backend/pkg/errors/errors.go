package errors

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AppError struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

func NewAppError(code int, msg string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: msg,
		Err:     err,
	}
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		for _, e := range c.Errors {
			if appErr, ok := e.Err.(*AppError); ok {
				c.JSON(appErr.Code, appErr)
				return
			}
		}
	}
}

// Auth & Session
var (
	ErrNotAuthed          = NewAppError(http.StatusUnauthorized, "Not authenticated", errors.New("not authenticated"))
	ErrInvalidToken       = NewAppError(http.StatusUnauthorized, "Invalid token", errors.New("invalid token"))
	ErrTokenExpired       = NewAppError(http.StatusUnauthorized, "Token has expired", errors.New("token has expired"))
	ErrInvalidCredentials = NewAppError(http.StatusUnauthorized, "Invalid email or password", errors.New("invalid credentials"))
	ErrPermissionDenied   = NewAppError(http.StatusForbidden, "Permission denied", errors.New("permission denied"))
	ErrTokenNotFound      = NewAppError(http.StatusUnauthorized, "Token not found in request", errors.New("token not found"))
	ErrMalformedToken     = NewAppError(http.StatusUnauthorized, "Malformed token", errors.New("malformed token"))
)

// User
var (
	ErrUserNotFound      = NewAppError(http.StatusNotFound, "User not found", errors.New("user not found"))
	ErrUserAlreadyExists = NewAppError(http.StatusConflict, "User already exists", errors.New("user already exists"))
	ErrInvalidUserID     = NewAppError(http.StatusBadRequest, "Invalid user ID", errors.New("invalid user ID"))
	ErrInvalidUserData   = NewAppError(http.StatusBadRequest, "Invalid user data", errors.New("invalid user data"))
)

// Validation
var (
	ErrValidationFailed     = NewAppError(http.StatusBadRequest, "Validation failed", errors.New("validation failed"))
	ErrMissingRequiredField = NewAppError(http.StatusBadRequest, "Missing required field", errors.New("missing required field"))
	ErrInvalidEmailFormat   = NewAppError(http.StatusBadRequest, "Invalid email format", errors.New("invalid email format"))
	ErrWeakPassword         = NewAppError(http.StatusBadRequest, "Weak password", errors.New("weak password"))
)

// Database
var (
	ErrDBConnection = NewAppError(http.StatusInternalServerError, "Database connection failed", errors.New("db connection failed"))
	ErrDBQuery      = NewAppError(http.StatusInternalServerError, "Database query error", errors.New("db query error"))
	ErrDBInsert     = NewAppError(http.StatusInternalServerError, "Insert error", errors.New("insert error"))
	ErrDBUpdate     = NewAppError(http.StatusInternalServerError, "Update error", errors.New("update error"))
	ErrDBDelete     = NewAppError(http.StatusInternalServerError, "Delete error", errors.New("delete error"))
)

// Internal / Server
var (
	ErrInternalServer     = NewAppError(http.StatusInternalServerError, "Internal server error", errors.New("internal error"))
	ErrServiceUnavailable = NewAppError(http.StatusServiceUnavailable, "Service unavailable", errors.New("service unavailable"))
	ErrUnexpected         = NewAppError(http.StatusInternalServerError, "Unexpected error occurred", errors.New("unexpected error"))
)

// Request
var (
	ErrBadRequest      = NewAppError(http.StatusBadRequest, "Bad request", errors.New("bad request"))
	ErrUnauthorized    = NewAppError(http.StatusUnauthorized, "Unauthorized", errors.New("unauthorized"))
	ErrForbidden       = NewAppError(http.StatusForbidden, "Forbidden", errors.New("forbidden"))
	ErrNotFound        = NewAppError(http.StatusNotFound, "Not found", errors.New("not found"))
	ErrConflict        = NewAppError(http.StatusConflict, "Conflict", errors.New("conflict"))
	ErrTooManyRequests = NewAppError(http.StatusTooManyRequests, "Too many requests", errors.New("too many requests"))
)

// File
var (
	ErrFileNotFound    = NewAppError(http.StatusNotFound, "File not found", errors.New("file not found"))
	ErrFileTooLarge    = NewAppError(http.StatusRequestEntityTooLarge, "File too large", errors.New("file too large"))
	ErrInvalidFileType = NewAppError(http.StatusUnsupportedMediaType, "Invalid file type", errors.New("invalid file type"))
)
