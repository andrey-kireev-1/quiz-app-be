package service

import (
	"errors"
	"quiz-app-be/internal/model"
	modeldb "quiz-app-be/internal/model/modelDB"
	modelhttp "quiz-app-be/internal/model/modelHttp"
	"quiz-app-be/internal/repository"
	"strings"
)

type UserService struct {
	repo *repository.Users
}

func NewUserService(repo *repository.Users) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(registerReq modelhttp.RegisterRequest) (modelhttp.TokenResponse, error) {
	if registerReq.Email == "" ||
		registerReq.Password == "" ||
		registerReq.Name == "" {
		return modelhttp.TokenResponse{}, errors.New(model.ErrEmptyFields)
	}

	// Check if user already exists
	_, err := s.repo.GetUserByEmail(registerReq.Email)
	if err != nil && strings.Contains(err.Error(), "no rows in result set") {
		err = nil
	}
	if err != nil && strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
		return modelhttp.TokenResponse{}, errors.New(model.ErrUserAlreadyExists)
	}
	if err != nil {
		return modelhttp.TokenResponse{}, err
	}

	hashedPassword, err := HashPassword(registerReq.Password)
	if err != nil {
		return modelhttp.TokenResponse{}, err
	}
	user := modeldb.User{
		Name:     registerReq.Name,
		Email:    registerReq.Email,
		Password: hashedPassword,
	}
	err = s.repo.CreateUser(user)
	if err != nil {
		return modelhttp.TokenResponse{}, err
	}

	accessToken, refreshToken, err := GenerateTokens(user.ID)
	if err != nil {
		return modelhttp.TokenResponse{}, errors.New(model.ErrGenerateTokens)
	}

	tokens := modelhttp.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return tokens, nil
}

func (s *UserService) LoginUser(loginReq modelhttp.LoginRequest) (modelhttp.TokenResponse, error) {
	if loginReq.Email == "" || loginReq.Password == "" {
		return modelhttp.TokenResponse{}, errors.New(model.ErrEmptyFields)
	}
	user, err := s.repo.GetUserByEmail(loginReq.Email)
	if err != nil && strings.Contains(err.Error(), "no rows in result set") {
		return modelhttp.TokenResponse{}, errors.New(model.ErrUserNotFound)
	}
	if err != nil {
		return modelhttp.TokenResponse{}, err
	}

	// Check password
	if !CheckPassword(loginReq.Password, user.Password) {
		return modelhttp.TokenResponse{}, errors.New(model.ErrInvalidPassword)
	}

	accessToken, refreshToken, err := GenerateTokens(user.ID)
	if err != nil {
		return modelhttp.TokenResponse{}, errors.New(model.ErrGenerateTokens)
	}
	// Create mock tokens
	tokens := modelhttp.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return tokens, nil
}

func (s *UserService) RefreshToken(refreshReq modelhttp.RefreshRequest) (modelhttp.TokenResponse, error) {
	userID, err := ValidateRefreshToken(refreshReq.RefreshToken)
	if err != nil {
		return modelhttp.TokenResponse{}, errors.New(model.ErrInvalidRefreshToken)
	}

	accessToken, refreshToken, err := GenerateTokens(userID)
	if err != nil {
		return modelhttp.TokenResponse{}, errors.New(model.ErrGenerateTokens)
	}

	tokens := modelhttp.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return tokens, nil
}

func (s *UserService) GetProfile(accessToken string) (modelhttp.ProfileResponse, error) {
	userID, err := ValidateAccessToken(accessToken)
	if err != nil {
		return modelhttp.ProfileResponse{}, err
	}

	user, err := s.repo.GetUserProfile(userID)
	if err != nil {
		return modelhttp.ProfileResponse{}, err
	}

	return modelhttp.ProfileResponse{
		Name:  user.Name,
		Email: user.Email,
	}, nil
}
