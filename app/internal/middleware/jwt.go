package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var secretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

type JWTClaims struct {
	GoogleClaims string `json:"googleClaims"`
	jwt.StandardClaims
}

func (m *middleware) GenerateJWTToken(googleClaims string, expiration time.Duration) (string, error) {
	claims := JWTClaims{
		GoogleClaims: googleClaims,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expiration).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func (m *middleware) ExtractJWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.String(http.StatusUnauthorized, "Unauthorized")
		}

		oauthProvider := c.Request().Header.Get("Oauth-Provider")
		if oauthProvider == "" {
			return c.String(http.StatusUnauthorized, "Unauthorized")
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return c.String(http.StatusUnauthorized, "Invalid token format")
		}

		if oauthProvider == "google" {
			accessTokenDetail, err := m.getAccessTokenDetailGoogle(tokenParts[1])
			if err != nil {
				return c.String(http.StatusUnauthorized, "Cannot validate token")
			}

			if accessTokenDetail.Error.Code == 401 {
				return c.String(http.StatusUnauthorized, "Invalid token")
			} else {
				c.Set("email", accessTokenDetail.Email)
				return next(c)
			}
		} else if oauthProvider == "github" {
			accessTokenDetail, err := m.getAccessTokenDetailGithub(tokenParts[1])
			if err != nil {
				fmt.Println(err)
				return c.String(http.StatusUnauthorized, "Cannot validate token")
			}
			c.Set("username", accessTokenDetail.User.Login)
			return next(c)
		}

		return nil

	}
}

func (m *middleware) getAccessTokenDetailGoogle(accessToken string) (AccessTokenDetail, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + url.QueryEscape(accessToken))
	if err != nil {
		return AccessTokenDetail{}, err
	}
	defer resp.Body.Close()

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return AccessTokenDetail{}, err
	}

	var accessTokenDetail AccessTokenDetail
	if err := json.Unmarshal(response, &accessTokenDetail); err != nil {
		return AccessTokenDetail{}, errors.Wrap(err, "failed to unmarshal JSON response")
	}

	return accessTokenDetail, nil
}

func (m *middleware) getAccessTokenDetailGithub(accessToken string) (GithubDetail, error) {

	clientID := viper.GetString("github.clientID")
	clientSecret := viper.GetString("github.clientSecret")

	url := fmt.Sprintf("https://api.github.com/applications/%s/token", clientID)
	payload := map[string]string{"access_token": accessToken}
	jsonPayload, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return GithubDetail{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/vnd.github+json")
	req.SetBasicAuth(clientID, clientSecret)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return GithubDetail{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return GithubDetail{}, fmt.Errorf("status: %v", resp.Status)
	}

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return GithubDetail{}, err
	}

	var githubDetail GithubDetail
	if err := json.Unmarshal(response, &githubDetail); err != nil {
		return GithubDetail{}, errors.Wrap(err, "failed to unmarshal JSON response")
	}

	fmt.Println("Response Status:", resp.Status)

	return githubDetail, nil
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
