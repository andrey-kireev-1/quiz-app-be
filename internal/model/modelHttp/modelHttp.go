package modelhttp

import (
	"encoding/json"
	"time"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateTestRequest struct {
	Title        string          `json:"title"`
	Description  string          `json:"description"`
	Image        *ImageData      `json:"image,omitempty"`
	IsStrict     bool            `json:"isStrict"`
	IsPrivate    bool            `json:"isPrivate"`
	Questions    json.RawMessage `json:"questions"`
	RefreshToken string          `json:"refreshToken"`
}

type ImageData struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Data string `json:"data"`
}

type GetTestResponse struct {
	Id          string `json:"id,omitempty"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"image,omitempty"`
	IsStrict    bool   `json:"isStrict"`
	IsPrivate   bool   `json:"isPrivate"`
	Questions   string `json:"questions"`
	AuthorName  string `json:"authorName"`
	CreatedAt   string `json:"createdAt"`
}

type SetResultRequest struct {
	TestID       string          `json:"testId"`
	Result       json.RawMessage `json:"result"`
	Score        int             `json:"score"`
	RefreshToken string          `json:"refreshToken"`
}

type GetAllTestsResponse struct {
	Tests []GetTestResponse `json:"tests"`
}

type ProfileResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type TestResultResponse struct {
	ID        string    `json:"id"`
	TestID    string    `json:"testId"`
	TestName  string    `json:"testName"`
	Score     int       `json:"score"`
	Result    string    `json:"result"`
	CreatedAt time.Time `json:"createdAt"`
}

type MyResultsResponse struct {
	Results []TestResultResponse `json:"results"`
}

type TestWithResultsResponse struct {
	TestID      string    `json:"testId"`
	TestName    string    `json:"testName"`
	UserName    string    `json:"userName"`
	Score       int       `json:"score"`
	Result      string    `json:"result"`
	CompletedAt time.Time `json:"completedAt"`
}

type MyTestsResultsResponse struct {
	Results []TestWithResultsResponse `json:"results"`
}
