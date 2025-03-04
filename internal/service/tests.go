package service

import (
	"encoding/base64"
	"quiz-app-be/internal/model"
	modeldb "quiz-app-be/internal/model/modelDB"
	modelhttp "quiz-app-be/internal/model/modelHttp"
	"quiz-app-be/internal/repository"
	"quiz-app-be/internal/setup/aws"
	"strconv"
	"strings"
	"time"

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

const testsLimit = 5

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

	var description, imageURL string
	if test.Description != nil {
		description = *test.Description
	}
	if test.ImageURL != nil {
		imageURL = *test.ImageURL
	}
	return modelhttp.GetTestResponse{
		Title:       test.Name,
		Description: description,
		ImageURL:    imageURL,
		IsStrict:    test.IsStrict,
		IsPrivate:   test.IsPrivate,
		Questions:   test.Data,
		AuthorName:  test.AuthorName,
		CreatedAt:   test.CreatedAt.String(),
	}, nil

}

func (s *TestService) GetHomeTests(paginationNumStr string) ([]modelhttp.GetTestResponse, error) {
	paginationNum, err := strconv.Atoi(paginationNumStr)
	if err != nil {
		return nil, err
	}
	if paginationNum < 1 {
		paginationNum = 1
	}
	limit := paginationNum * testsLimit
	offset := limit - testsLimit
	limit = limit - offset
	tests, err := s.repo.GetHomeTests(limit, offset)
	if err != nil {
		return nil, err
	}

	var resp []modelhttp.GetTestResponse
	for _, test := range tests {
		var description, imageURL string
		if test.Description != nil {
			description = *test.Description
		}
		if test.ImageURL != nil {
			imageURL = *test.ImageURL
		}
		resp = append(resp, modelhttp.GetTestResponse{
			Id:          test.ID.String(),
			Title:       test.Name,
			Description: description,
			ImageURL:    imageURL,
			IsStrict:    test.IsStrict,
			IsPrivate:   test.IsPrivate,
			Questions:   test.Data,
			AuthorName:  test.AuthorName,
			CreatedAt:   test.CreatedAt.String(),
		})
	}
	return resp, nil
}

func (s *TestService) CountAllPublicTests() (int, error) {
	return s.repo.CountAllPublicTests()
}

func (s *TestService) GetUserTests(accessToken string) (modelhttp.MyTestsResponse, error) {
	authorID, err := ValidateAccessToken(accessToken)
	if err != nil {
		return modelhttp.MyTestsResponse{}, err
	}

	tests, err := s.repo.GetUserTests(authorID)
	if err != nil {
		return modelhttp.MyTestsResponse{}, err
	}

	response := modelhttp.MyTestsResponse{
		Tests: make([]modelhttp.GetTestResponse, len(tests)),
	}

	for i, test := range tests {
		var description, imageURL string
		if test.Description != nil {
			description = *test.Description
		}
		if test.ImageURL != nil {
			imageURL = *test.ImageURL
		}
		response.Tests[i] = modelhttp.GetTestResponse{
			Id:          test.ID.String(),
			Title:       test.Name,
			Description: description,
			ImageURL:    imageURL,
			IsStrict:    test.IsStrict,
			IsPrivate:   test.IsPrivate,
			Questions:   test.Data,
			CreatedAt:   test.CreatedAt.Format(time.RFC3339),
		}
	}

	return response, nil
}

func (s *TestService) GetFilteredTests(filters modelhttp.TestFilters) (modelhttp.GetAllTestsResponse, error) {
	tests, err := s.repo.GetFilteredTests(filters)
	if err != nil {
		return modelhttp.GetAllTestsResponse{}, err
	}

	response := modelhttp.GetAllTestsResponse{
		Tests: make([]modelhttp.GetTestResponse, len(tests)),
	}

	for i, test := range tests {
		var description, imageURL string
		if test.Description != nil {
			description = *test.Description
		}
		if test.ImageURL != nil {
			imageURL = *test.ImageURL
		}
		response.Tests[i] = modelhttp.GetTestResponse{
			Id:          test.ID.String(),
			Title:       test.Name,
			Description: description,
			ImageURL:    imageURL,
			IsStrict:    test.IsStrict,
			IsPrivate:   test.IsPrivate,
			Questions:   test.Data,
			AuthorName:  test.AuthorName,
			CreatedAt:   test.CreatedAt.Format(time.RFC3339),
		}
	}

	return response, nil
}
