package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gamepkw/google-oauth2-user-service/app/internal/middleware"
	"github.com/gamepkw/google-oauth2-user-service/app/internal/model"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type UserHandler struct {
}

func NewUserHandler(e *echo.Echo) {
	handler := &UserHandler{}

	middL := middleware.InitMiddleware()
	resourceGroup := e.Group("/data", middL.ExtractJWTMiddleware)
	resourceGroup.GET("/get-my-email", handler.GetEmail)
	resourceGroup.GET("/get-my-username", handler.GetUsername)
	resourceGroup.GET("/get-secret-content", handler.GetSecretContent)
	resourceGroup.GET("/check-token", handler.CheckToken)
}

func (a *UserHandler) GetEmail(c echo.Context) error {
	accessTokenDetail := c.Get("accessTokenDetail").([]byte)

	var response model.AccessTokenDetail
	if err := json.Unmarshal(accessTokenDetail, &response); err != nil {
		return errors.Wrap(err, "failed to unmarshal JSON response")
	}

	fmt.Println(response.Email)

	return c.JSON(http.StatusOK, response.Email)
}

func (a *UserHandler) GetUsername(c echo.Context) error {

	accessTokenDetail := c.Get("accessTokenDetail").([]byte)

	var response model.GithubDetail
	if err := json.Unmarshal(accessTokenDetail, &response); err != nil {
		return errors.Wrap(err, "failed to unmarshal JSON response")
	}

	fmt.Println(response.User.Login)

	return c.JSON(http.StatusOK, response.User.Login)
}

func (a *UserHandler) CheckToken(c echo.Context) error {
	accessTokenDetail := c.Get("accessTokenDetail").([]byte)

	fmt.Println(string(accessTokenDetail))

	return c.JSON(http.StatusOK, string(accessTokenDetail))
}

func (a *UserHandler) GetSecretContent(c echo.Context) error {
	// content := "This content is show only for authorized users"

	email := c.Get("email").(string)

	var name string
	for _, user := range users {
		if user.Email == email {
			name = user.Name
			break
		}
	}

	resp := "Your name is " + name

	return c.JSON(http.StatusOK, resp)
}

type User struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

var users = []User{
	{Email: "game9074@gmail.com", Name: "Pakawat Bamrungkit"},
	{Email: "gamepkw9074@gmail.com", Name: "Game Pakawat"},
}
