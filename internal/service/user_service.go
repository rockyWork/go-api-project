package service

import (
	"context"
	"errors"

	"go-api-project/internal/model"
	"go-api-project/internal/model/dto"
	"go-api-project/internal/repository"
	"go-api-project/pkg/jwt"
	"go-api-project/pkg/utils"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUsernameExists     = errors.New("username already exists")
	ErrEmailExists        = errors.New("email already exists")
	ErrUserBanned         = errors.New("user has been banned")
	ErrInvalidPassword    = errors.New("invalid password")
	ErrSamePassword       = errors.New("new password cannot be the same as old password")
)

type UserService struct {
	userRepo    *repository.UserRepository
	jwtManager  *jwt.Manager
}

func NewUserService(userRepo *repository.UserRepository, jwtManager *jwt.Manager) *UserService {
	return &UserService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

func (s *UserService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.TokenResponse, error) {
	if s.userRepo.ExistsUsername(ctx, req.Username) {
		return nil, ErrUsernameExists
	}

	if s.userRepo.ExistsEmail(ctx, req.Email) {
		return nil, ErrEmailExists
	}

	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		Phone:        req.Phone,
		PasswordHash: passwordHash,
		Status:       model.UserStatusNormal,
		Role:         model.UserRoleUser,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	accessToken, refreshToken, err := s.jwtManager.GenerateTokenPair(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	return &dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    7200,
		TokenType:    "Bearer",
		User:         s.toUserResponse(user),
	}, nil
}

func (s *UserService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.TokenResponse, error) {
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if !user.IsNormal() {
		return nil, ErrUserBanned
	}

	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	accessToken, refreshToken, err := s.jwtManager.GenerateTokenPair(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	return &dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    7200,
		TokenType:    "Bearer",
		User:         s.toUserResponse(user),
	}, nil
}

func (s *UserService) RefreshToken(ctx context.Context, refreshToken string) (*dto.TokenResponse, error) {
	claims, err := s.jwtManager.ParseRefreshToken(refreshToken)
	if err != nil {
		if errors.Is(err, jwt.ErrExpiredToken) {
			return nil, errors.New("refresh token expired")
		}
		return nil, errors.New("invalid refresh token")
	}

	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	if !user.IsNormal() {
		return nil, ErrUserBanned
	}

	newAccessToken, newRefreshToken, err := s.jwtManager.RefreshTokens(refreshToken)
	if err != nil {
		return nil, err
	}

	return &dto.TokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    7200,
		TokenType:    "Bearer",
	}, nil
}

func (s *UserService) ChangePassword(ctx context.Context, userID uint, req *dto.ChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if !utils.CheckPassword(req.OldPassword, user.PasswordHash) {
		return ErrInvalidPassword
	}

	if utils.CheckPassword(req.NewPassword, user.PasswordHash) {
		return ErrSamePassword
	}

	newPasswordHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	user.PasswordHash = newPasswordHash
	return s.userRepo.Update(ctx, user)
}

func (s *UserService) GetUserByID(ctx context.Context, id uint) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.toUserResponse(user), nil
}

func (s *UserService) GetUserList(ctx context.Context, req *dto.UserListRequest) (*dto.UserListResponse, error) {
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	users, total, err := s.userRepo.List(ctx, req.Page, req.PageSize, req.Keyword, req.Status, req.Role)
	if err != nil {
		return nil, err
	}

	list := make([]*dto.UserResponse, len(users))
	for i, user := range users {
		list[i] = s.toUserResponse(user)
	}

	return &dto.UserListResponse{
		List:     list,
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    total,
	}, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id uint, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Email != "" && req.Email != user.Email {
		if s.userRepo.ExistsEmail(ctx, req.Email) {
			return nil, ErrEmailExists
		}
		user.Email = req.Email
	}

	if req.Phone != "" {
		user.Phone = req.Phone
	}

	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return s.toUserResponse(user), nil
}

func (s *UserService) DeleteUser(ctx context.Context, id uint) error {
	return s.userRepo.Delete(ctx, id)
}

func (s *UserService) toUserResponse(user *model.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Phone:     user.Phone,
		Avatar:    user.Avatar,
		Status:    user.Status,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}
}
