package app

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewServer(microService *CryptoMicroservice) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	v1 := e.Group("/v1")
	v1.POST("/challenge", microService.CreateChallenge)
	v1.POST("/verify-challenge", microService.VerifyChallenge)

	return e
}
