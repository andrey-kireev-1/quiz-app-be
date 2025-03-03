package service

import (
	"encoding/base64"
	"quiz-app-be/internal/model"
	modeldb "quiz-app-be/internal/model/modelDB"
	modelhttp "quiz-app-be/internal/model/modelHttp"
	"quiz-app-be/internal/repository"
	"quiz-app-be/internal/setup/aws"
	"strings"

	"github.com/pkg/errors"
)

type TestService struct {
	repo      *repository.Tests
	usersRepo *repository.Users

	s3 *aws.AwsClient
}

func NewTestService(
	repo *repository.Tests,
	usersRepo *repository.Users,
	s3 *aws.AwsClient,
) *TestService {
	return &TestService{
		repo:      repo,
		usersRepo: usersRepo,
		s3:        s3,
	}
}

func (s *TestService) CreateTest(req modelhttp.CreateTestRequest) error {
	userID, err := ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return err
	}

	exists, err := s.usersRepo.CheckUserByID(userID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New(model.ErrUserNotFound)
	}

	var url string
	if req.Image != nil {
		base64Data := req.Image.Data
		if idx := strings.Index(base64Data, ","); idx != -1 {
			base64Data = base64Data[idx+1:]
		}

		// Decode base64 to bytes
		imageBytes, err := base64.StdEncoding.DecodeString(base64Data)
		if err != nil {
			return err
		}

		err = s.s3.Upload(
			imageBytes,
			req.Image.Name,
			req.Image.Type,
		)
		if err != nil {
			return err
		}

		url, err = s.s3.GetDownloadUrlByName(req.Image.Name)
		if err != nil {
			return err
		}
	}

	test := modeldb.Test{
		Name:        req.Title,
		Description: &req.Description,
		Data:        string(req.Questions),
		ImageURL:    &url,
		IsStrict:    req.IsStrict,
		IsPrivate:   req.IsPrivate,
		AuthorID:    userID,
	}
	err = s.repo.CreateTest(test)
	if err != nil {
		return err
	}
	return nil
}

func (s *TestService) GetTest(testID string) (modelhttp.GetTestResponse, error) {
	test, err := s.repo.GetTestByID(testID)
	if err != nil {
		return modelhttp.GetTestResponse{}, err
	}

	return modelhttp.GetTestResponse{
		Title:       test.Name,
		Description: *test.Description,
		ImageURL:    *test.ImageURL,
		IsStrict:    test.IsStrict,
		IsPrivate:   test.IsPrivate,
		Questions:   test.Data,
		AuthorName:  test.AuthorID.String(),
		CreatedAt:   test.CreatedAt.String(),
	}, nil

}
