package auth

import (
	authpb "backend/generated/proto/auth"
	"backend/internal/database/postgres"
	"backend/pkg/config"
	"backend/pkg/security"
	jwt "backend/pkg/security/JWT"
	"backend/pkg/types"
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	authpb.UnimplementedAuthServiceServer

	UserDataBase   *postgres.PostgresDatabase
	cfg            *config.Config
	logger         *zap.Logger
	authController *jwt.AuthController
}

func NewAuthServer(
	db *postgres.PostgresDatabase,
	cfg *config.Config,
	logger *zap.Logger,
	authController *jwt.AuthController) *AuthServer {
	return &AuthServer{
		UserDataBase:   db,
		cfg:            cfg,
		logger:         logger,
		authController: authController,
	}
}
func (s *AuthServer) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.AuthResponse, error) {
	user, err := s.UserDataBase.GetUser(map[string]interface{}{"nickname": req.Nickname})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	err = security.CompareHashAndPassword(req.Password, user.HashPassword, s.cfg)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid password")
	}

	accessToken, err := s.authController.CreateAccessToken(types.UserInfo{Id: int64(user.ID), Email: user.Email, Nickname: user.Nickname})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate access token: %v", err)
	}

	refreshToken, err := s.authController.CreateRefreshToken(types.UserInfo{Id: int64(user.ID), Email: user.Email, Nickname: user.Nickname})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate refresh token: %v", err)
	}
	return &authpb.AuthResponse{
		Status:       int64(codes.OK),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Email:        user.Email,
	}, nil
}

// Register реализует AuthService.Register
func (s *AuthServer) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.AuthResponse, error) {
	// 1. Валидация входных данных
	if req.Nickname == "" || req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "nickname, email and password are required")
	}

	if len(req.Password) < 8 {
		return nil, status.Error(codes.InvalidArgument, "password must be at least 8 characters")
	}

	// 2. Проверка, что пользователь еще не существует
	existingUser, err := s.UserDataBase.GetUser(map[string]interface{}{"email": req.Email})
	if existingUser != nil && err == nil {
		return nil, status.Error(codes.AlreadyExists, "user with this email already exists")
	}

	existingUser, err = s.UserDataBase.GetUser(map[string]interface{}{"nickname": req.Nickname})
	if existingUser != nil && err == nil {
		return nil, status.Error(codes.AlreadyExists, "user with this nickname already exists")
	}

	// 3. Хеширование пароля
	hashedPassword, err := security.GenerateHashedPassword(req.Password)
	if err != nil {
		s.logger.Error("Failed to hash password", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to process password")
	}

	// 4. Создание пользователя
	newUser := &postgres.User{
		Nickname:     req.Nickname,
		Email:        req.Email,
		HashPassword: hashedPassword,
	}

	if err := s.UserDataBase.AddUser(newUser); err != nil {
		s.logger.Error("Failed to create user", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to create user")
	}

	// 5. Генерация токенов
	accessToken, err := s.authController.CreateAccessToken(types.UserInfo{Id: int64(newUser.ID), Email: newUser.Email, Nickname: newUser.Nickname})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate access token: %v", err)
	}

	refreshToken, err := s.authController.CreateRefreshToken(types.UserInfo{Id: int64(newUser.ID), Email: newUser.Email, Nickname: newUser.Nickname})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate refresh token: %v", err)
	}

	// 6. Возвращаем ответ
	return &authpb.AuthResponse{
		Status:       int64(codes.OK),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Email:        newUser.Email,
	}, nil
}

// Refresh реализует AuthService.Refresh
func (s *AuthServer) Refresh(ctx context.Context, req *authpb.RefreshRequest) (*authpb.AuthResponse, error) {
	// 1. Валидация входного токена
	if req.RefreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh token is required")
	}

	// 2. Парсинг и валидация refresh token
	claims, err := s.authController.Decrypt(req.RefreshToken)
	if err != nil {
		s.logger.Debug("Invalid refresh token", zap.Error(err))
		return nil, status.Error(codes.Unauthenticated, "invalid refresh token")
	}

	if tokenType, ok := claims["type"].(string); !ok || tokenType != "refresh" {
		return nil, status.Error(codes.Unauthenticated, "not a refresh token")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, status.Error(codes.Internal, "invalid token claims")
	}

	user, err := s.UserDataBase.GetUser(map[string]interface{}{"nickname": email})
	if err != nil {
		s.logger.Error("User not found", zap.String("email", email), zap.Error(err))
		return nil, status.Error(codes.NotFound, "user not found")
	}

	// 5. Генерация новых токенов
	accessToken, err := s.authController.CreateAccessToken(types.UserInfo{Id: int64(user.ID), Email: user.Email, Nickname: user.Nickname})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate access token: %v", err)
	}

	refreshToken, err := s.authController.CreateRefreshToken(types.UserInfo{Id: int64(user.ID), Email: user.Email, Nickname: user.Nickname})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate refresh token: %v", err)
	}

	// 6. Возвращаем новые токены
	return &authpb.AuthResponse{
		Status:       int64(codes.OK),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Email:        user.Email,
	}, nil
}

// Me реализует AuthService.Me
func (s *AuthServer) Me(ctx context.Context, req *authpb.GetMeRequest) (*authpb.AuthResponse, error) {
	// 2. Парсим и валидируем токен
	claims, err := s.authController.Decrypt(req.RefreshToken)
	if err != nil {
		s.logger.Debug("Invalid refresh token", zap.Error(err))
		return nil, status.Error(codes.Unauthenticated, "invalid refresh token")
	}

	if tokenType, ok := claims["type"].(string); !ok || tokenType != "refresh" {
		return nil, status.Error(codes.Unauthenticated, "not a refresh token")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, status.Error(codes.Internal, "invalid token claims")
	}

	user, err := s.UserDataBase.GetUser(map[string]interface{}{"nickname": email})
	if err != nil {
		s.logger.Error("User not found", zap.String("email", email), zap.Error(err))
		return nil, status.Error(codes.NotFound, "user not found")
	}
	accessToken, err := s.authController.CreateAccessToken(types.UserInfo{Id: int64(user.ID), Email: user.Email, Nickname: user.Nickname})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate access token: %v", err)
	}

	refreshToken, err := s.authController.CreateRefreshToken(types.UserInfo{Id: int64(user.ID), Email: user.Email, Nickname: user.Nickname})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate refresh token: %v", err)
	}

	return &authpb.AuthResponse{
		Status:       int64(codes.OK),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Email:        user.Email,
	}, nil
}
