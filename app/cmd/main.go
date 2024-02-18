package main

import (
	"log"
	"net/http"

	"github.com/gamepkw/google-oauth2-user-service/app/internal/config"
	_userHandler "github.com/gamepkw/google-oauth2-user-service/app/internal/handler"
	"github.com/gamepkw/google-oauth2-user-service/tools/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
)

func main() {
	config.InitializeViper()
	logger.InitializeZapCustomLogger()

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3001", "http://localhost:3002"},
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		AllowCredentials: true,
	}))

	_userHandler.NewUserHandler(e)
	logger.Log.Info("Started running on http://localhost:" + viper.GetString("port"))
	log.Fatal(e.Start(":" + viper.GetString("port")))
}
