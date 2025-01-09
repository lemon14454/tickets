package request

import (
	"fmt"
	"net/http"
	"time"
)

type userResponse struct {
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Host      bool      `json:"bool"`
	CreatedAt time.Time `json:"created_at"`
}

type loginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginUserResponse struct {
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  userResponse `json:"user"`
}

func (client *Client) LoginUser(username, password string) (*loginUserResponse, error) {
	user, err := MakeRequest[loginUserResponse](client, http.MethodPost, "users/login", loginUserRequest{
		Email:    fmt.Sprintf("%s@mail.com", username),
		Password: password,
	}, nil)
	return user, err
}

type registerUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Host     bool   `json:"host"`
}

func (client *Client) RegisterUser(username, mail, password string) (*userResponse, error) {
	user, err := MakeRequest[userResponse](client, http.MethodPost, "users", registerUserRequest{
		Username: username,
		Email:    mail,
		Password: password,
		Host:     true,
	}, nil)
	return user, err
}

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type renewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (client *Client) RenewToken(refreshToken string) (*renewAccessTokenResponse, error) {
	res, err := MakeRequest[renewAccessTokenResponse](client, http.MethodPost, "tokens/renew", renewAccessTokenRequest{
		RefreshToken: refreshToken,
	}, nil)
	return res, err
}
