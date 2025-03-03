package modelhttp

import "encoding/json"

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

type GetTestRequest struct {
	ID string `json:"id"`
}

type GetTestResponse struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"image,omitempty"`
	IsStrict    bool   `json:"isStrict"`
	IsPrivate   bool   `json:"isPrivate"`
	Questions   string `json:"questions"`
	AuthorName  string `json:"authorName"`
	CreatedAt   string `json:"createdAt"`
}
