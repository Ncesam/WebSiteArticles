package auth

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	authpb "backend/generated/proto/auth"
	"backend/internal/database/postgres"
	"backend/pkg/config"
	"backend/pkg/security"
	jwt "backend/pkg/security/JWT"
	"backend/pkg/types"
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
	logger.Info("Created Auth Server")
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
		return nil, status.Errorf(codes.NotFound, "failed to get user: %v", err)
	}
	s.logger.Debug("got user data", zap.Int64("user id", user.ID))

	err = security.CompareHashAndPassword(req.Password, user.HashPassword)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid password")
	}
	s.logger.Debug("User password is valid")

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

func (s *AuthServer) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.AuthResponse, error) {
	if req.Nickname == "" || req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "nickname, email and password are required")
	}
	s.logger.Debug("Got user data", zap.String("email", req.Email))
	if len(req.Password) < 8 {
		return nil, status.Error(codes.InvalidArgument, "password must be at least 8 characters")
	}

	existingUser, err := s.UserDataBase.GetUser(map[string]interface{}{"email": req.Email})
	if existingUser != nil && err == nil {
		return nil, status.Error(codes.AlreadyExists, "user with this email already exists")
	}

	existingUser, err = s.UserDataBase.GetUser(map[string]interface{}{"nickname": req.Nickname})
	if existingUser != nil && err == nil {
		return nil, status.Error(codes.AlreadyExists, "user with this nickname already exists")
	}

	hashedPassword, err := security.GenerateHashedPassword(req.Password)
	if err != nil {
		s.logger.Error("Failed to hash password", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to process password")
	}
	s.logger.Debug("User password is hashed")

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

func (s *AuthServer) Refresh(ctx context.Context, req *authpb.RefreshRequest) (*authpb.AuthResponse, error) {
	if req.RefreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh token is required")
	}
	s.logger.Debug("Got refresh token")
	claims, err := s.authController.DecryptRefresh(req.RefreshToken)
	if err != nil {
		s.logger.Debug("Invalid refresh token", zap.Error(err))
		return nil, status.Error(codes.Unauthenticated, "invalid refresh token")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, status.Error(codes.Internal, "invalid token claims")
	}
	s.logger.Debug("Got user email", zap.String("email", email))
	user, err := s.UserDataBase.GetUser(map[string]interface{}{"email": email})
	if err != nil {
		s.logger.Error("User not found", zap.String("email", email), zap.Error(err))
		return nil, status.Error(codes.NotFound, "user not found")
	}
	s.logger.Debug("Got user data", zap.Int64("user id", user.ID))
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

func (s *AuthServer) Me(ctx context.Context, req *authpb.GetMeRequest) (*authpb.AuthResponse, error) {
	s.logger.Debug("Recieve user data", zap.String("email", req.Email))
	user, err := s.UserDataBase.GetUser(map[string]interface{}{"email": req.Email})
	if err != nil {
		s.logger.Error("User not found", zap.String("email", req.Email), zap.Error(err))
		return nil, status.Error(codes.NotFound, "user not found")
	}
	s.logger.Debug("Got user", zap.Int64("user id", user.ID))
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
