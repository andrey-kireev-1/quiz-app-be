package service

import (
	"quiz-app-be/internal/model"
	modeldb "quiz-app-be/internal/model/modelDB"
	modelhttp "quiz-app-be/internal/model/modelHttp"
	"quiz-app-be/internal/repository"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type ResultService struct {
	repo      *repository.Results
	usersRepo *repository.Users
}

func NewResultService(
	repo *repository.Results,
	usersRepo *repository.Users,
) *ResultService {
	return &ResultService{
		repo:      repo,
		usersRepo: usersRepo,
	}
}

func (s *ResultService) SetResult(req modelhttp.SetResultRequest) error {
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

	result := modeldb.Result{
		UserID: userID,
		TestID: uuid.MustParse(req.TestID),
		Result: string(req.Result),
		Score:  req.Score,
	}
	return s.repo.SetResult(result)
}

func (s *ResultService) GetUserResults(accessToken string) (modelhttp.MyResultsResponse, error) {
	userID, err := ValidateAccessToken(accessToken)
	if err != nil {
		return modelhttp.MyResultsResponse{}, err
	}

	results, err := s.repo.GetUserResults(userID)
	if err != nil {
		return modelhttp.MyResultsResponse{}, err
	}

	response := modelhttp.MyResultsResponse{
		Results: make([]modelhttp.TestResultResponse, len(results)),
	}

	for i, result := range results {
		response.Results[i] = modelhttp.TestResultResponse{
			ID:        result.ID.String(),
			TestID:    result.TestID.String(),
			TestName:  result.TestName,
			Score:     result.Score,
			Result:    result.Result,
			CreatedAt: result.CreatedAt,
		}
	}

	return response, nil
}

func (s *ResultService) GetAuthorTestsResults(accessToken string) (modelhttp.MyTestsResultsResponse, error) {
	// Validate refresh token and get user ID
	authorID, err := ValidateAccessToken(accessToken)
	if err != nil {
		return modelhttp.MyTestsResultsResponse{}, err
	}

	// Get results from repository
	results, err := s.repo.GetResultsForAuthorTests(authorID)
	if err != nil {
		return modelhttp.MyTestsResultsResponse{}, err
	}

	// Convert to response model
	response := modelhttp.MyTestsResultsResponse{
		Results: make([]modelhttp.TestWithResultsResponse, len(results)),
	}

	for i, result := range results {
		response.Results[i] = modelhttp.TestWithResultsResponse{
			TestID:      result.TestID.String(),
			TestName:    result.TestName,
			UserName:    result.UserName,
			Score:       result.Score,
			Result:      result.Result,
			CompletedAt: result.CreatedAt,
		}
	}

	return response, nil
}
