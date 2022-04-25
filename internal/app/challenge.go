package app

import (
	"crypto-project-1/internal/domain"
	"crypto-project-1/public"
	//"encoding/hex"
	"encoding/json"
	"github.com/labstack/echo/v4"
	logger "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

// POST v1/challenge
func (m *CryptoMicroservice) CreateChallenge(ctx echo.Context) error {
	logger.Info("create challenge request received")

	requestBody, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		logger.Error(domain.CryptoAPIError, domain.UnexpectedError, "could not read request ", err)
		return ctx.JSON(http.StatusBadRequest, public.ApiResponse{
			Code:    public.ChallengeCreateFailed,
			Message: "could not read request body",
		})
	}
	defer func() {
		if err := ctx.Request().Body.Close(); err != nil {
			logger.Error(domain.CryptoAPIError, domain.UnexpectedError, "could not close request body ", err)
			return
		}
	}()
	request := &public.CreateChallengeRequestBody{}
	if err := json.Unmarshal(requestBody, request); err != nil {
		logger.Error(domain.CryptoAPIError, domain.UnexpectedError,
			"could not unmarshal request ", string(requestBody), err)
		return ctx.JSON(http.StatusBadRequest, public.ApiResponse{
			Code:    public.ChallengeCreateFailed,
			Message: "invalid request body",
		})
	}

	challenge, err := m.challengeService.CreateChallenge(request.PubKey)
	if err != nil {
		logger.Error(domain.CryptoAPIError, domain.UnexpectedError, "failed to create challenge ", err)
		return ctx.JSON(http.StatusInternalServerError, public.ApiResponse{
			Code:    public.ChallengeCreateFailed,
			Message: "error while trying to create challenge",
		})
	}

	logger.Info("challenge created successfully")
	return ctx.JSON(http.StatusOK, public.ApiResponse{
		Result:  challenge,
		Code:    public.ChallengeCreateSucceed,
		Message: "successfully created challenge",
	})
}

// POST v1/verify-challenge
func (m *CryptoMicroservice) VerifyChallenge(ctx echo.Context) error {
	logger.Info("verify challenge request received")

	requestBody, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		logger.Error(domain.CryptoAPIError, domain.UnexpectedError, "could not read request ", err)
		return ctx.JSON(http.StatusBadRequest, public.ApiResponse{
			Result: &domain.ChallengeValidationResult{
				Valid: false,
			},
			Code:    public.ChallengeValidationFailed,
			Message: "could not read request body",
		})
	}
	defer func() {
		if err := ctx.Request().Body.Close(); err != nil {
			logger.Error(domain.CryptoAPIError, domain.UnexpectedError, "could not close request body ", err)
			return
		}
	}()

	request := &public.VerifyChallengeRequestBody{}
	if err := json.Unmarshal(requestBody, request); err != nil {
		logger.Error(domain.CryptoAPIError, domain.UnexpectedError, "could not unmarshal request ", err)
		return ctx.JSON(http.StatusBadRequest, public.ApiResponse{
			Result: &domain.ChallengeValidationResult{
				Valid: false,
			},
			Code:    public.ChallengeValidationFailed,
			Message: "invalid request body",
		})
	}

	result, err := m.challengeService.VerifyChallenge(request.Token)
	if err != nil {
		logger.Error(domain.CryptoAPIError, domain.UnexpectedError, "error while challenge validation ", err)
		return ctx.JSON(http.StatusInternalServerError, public.ApiResponse{
			Result: &domain.ChallengeValidationResult{
				Valid: false,
			},
			Code:    public.ChallengeValidationFailed,
			Message: "internal error while trying to validate challenge",
		})
	}

	if !result.Valid {
		message := "challenge validation failed"
		logger.Info(message)
		return ctx.JSON(http.StatusOK, public.ApiResponse{
			Result:  result,
			Code:    public.ChallengeValidationFailed,
			Message: message,
		})
	}

	message := "challenge validation succeeded"
	logger.Info(message)
	return ctx.JSON(http.StatusOK, public.ApiResponse{
		Result:  result,
		Code:    public.ChallengeValidationsucceeded,
		Message: message,
	})
}
