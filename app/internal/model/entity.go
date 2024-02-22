package model

type CallbackGoogleRequest struct {
	AuthCode string `json:"authCode"`
}

type CallbackGoogleResponse struct {
	AccessToken string `json:"accessToken"`
}

type AccessTokenDetail struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Error         Error  `json:"error,omitempty"`
}

type GithubDetail struct {
	User GithubDetailUser `json:"user"`
}

type GithubDetailUser struct {
	Login string `json:"login"`
	ID    int    `json:"id"`
}

type Error struct {
	Code int `json:"code"`
}
