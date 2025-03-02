package modelhttp

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
