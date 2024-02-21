package model

type CallbackGoogleRequest struct {
	AuthCode string `json:"authCode"`
}

type CallbackGoogleResponse struct {
	AccessToken string `json:"accessToken"`
}
