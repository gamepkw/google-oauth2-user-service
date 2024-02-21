package handler

import (
	"fmt"
	"net/http"

	"github.com/gamepkw/google-oauth2-user-service/app/internal/middleware"
	"github.com/labstack/echo/v4"
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
}

func (a *UserHandler) GetEmail(c echo.Context) error {
	email := c.Get("email").(string)

	fmt.Println(email)

	return c.JSON(http.StatusOK, email)
}

func (a *UserHandler) GetUsername(c echo.Context) error {
	username := c.Get("username").(string)

	fmt.Println(username)

	return c.JSON(http.StatusOK, username)
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
