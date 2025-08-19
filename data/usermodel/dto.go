package usermodel

import "time"

type LoginDTO struct {
	SessionId                  string    `json:"session_id"`
	AccessToken                string    `json:"access_token"`
	AccessTokenExpirationTime  time.Time `json:"access_token_expiration"`
	RefreshToken               string    `json:"refresh_token"`
	RefreshTokenExpirationTime time.Time `json:"refresh_token_expiration_time"`
}

type RenewAccessToken struct {
	AccessToken               string    `json:"access_token"`
	AccessTokenExpirationTime time.Time `json:"access_token_expiration"`
}

type AccessTokenDTO struct {
	AccessToken string `json:"access_token"`
}
