package main

import (
	"crypto-project-1/internal/app"
	"crypto-project-1/internal/domain"
	"crypto-project-1/internal/repository"
	"crypto-project-1/internal/service"
	logger "github.com/sirupsen/logrus"
)

const (
	port = "7777"
)

func main() {
	logger.Info("crypto api starting...")

	db, err := repository.NewDB()
	if err != nil {
		logger.Error(domain.CryptoAPIError, domain.BootError, "could not connect to db ", err)
		return
	}
	err = db.Ping()
	if err != nil {
		logger.Error(domain.CryptoAPIError, domain.BootError, "cannot ping db ", err)
	}

	// initialize dependencies
	repository := repository.NewRepo()
	microservice := app.NewCryptoMicroservice(service.NewChallengeService(repository))

	// create routes
	httpServer := app.NewServer(microservice)
	// start http server
	if err := httpServer.Start(":" + port); err != nil {
		logger.Error(domain.CryptoAPIError, domain.BootError, "cannot start http server ", err)
		panic(err)
	}
}
