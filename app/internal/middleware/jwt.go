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
			c.Set("accessTokenDetail", accessTokenDetail)
			return next(c)
		} else if oauthProvider == "github" {
			accessTokenDetail, err := m.getAccessTokenDetailGithub(tokenParts[1])
			if err != nil {
				return c.String(http.StatusUnauthorized, "Cannot validate token")
			}
			c.Set("accessTokenDetail", accessTokenDetail)
			return next(c)
		}

		return nil

	}
}

func (m *middleware) getAccessTokenDetailGoogle(accessToken string) ([]byte, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + url.QueryEscape(accessToken))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, errors.New("Unauthorized")
	} else if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Error requesting access token")
	}

	defer resp.Body.Close()

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println("Response:", string(response))

	return response, nil
}

func (m *middleware) getAccessTokenDetailGithub(accessToken string) ([]byte, error) {

	clientID := viper.GetString("github.clientID")
	clientSecret := viper.GetString("github.clientSecret")

	url := fmt.Sprintf("https://api.github.com/applications/%s/token", clientID)
	payload := map[string]string{"access_token": accessToken}
	jsonPayload, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/vnd.github+json")
	req.SetBasicAuth(clientID, clientSecret)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, errors.New("Unauthorized")
	} else if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Error requesting access token")
	}

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	formattedResponse := formatJSON(response)

	fmt.Println("Response:", string(formattedResponse))

	return formattedResponse, nil
}

func formatJSON(jsonBytes []byte) []byte {
	var data interface{}
	err := json.Unmarshal(jsonBytes, &data)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return nil
	}

	formattedJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return nil
	}

	return []byte(formattedJSON)
}
