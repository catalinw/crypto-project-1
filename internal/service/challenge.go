package service

import (
	"bytes"
	"compress/gzip"
	"crypto-project-1/internal/domain"
	"crypto-project-1/internal/repository"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	logger "github.com/sirupsen/logrus"
	"io/ioutil"
	"time"
)

type ChallengeService interface {
	CreateChallenge(string) (*domain.Challenge, error)
	VerifyChallenge(string) (*domain.ChallengeValidationResult, error)
}

type challengeService struct {
	repo *repository.Repo
}

const (
	nonceTimeToLive = time.Minute * 5
)

func NewChallengeService(repo *repository.Repo) ChallengeService {
	return &challengeService{
		repo: repo,
	}
}

func (cs *challengeService) CreateChallenge(pubKey string) (*domain.Challenge, error) {
	return cs.repo.ChallengeRepo.CreateChallenge(pubKey, uuid.NewString(), time.Now().Add(nonceTimeToLive).Unix())
}

func (cs *challengeService) VerifyChallenge(signedToken string) (*domain.ChallengeValidationResult, error) {
	claims := &jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(signedToken, claims, getPublicKey)
	if err != nil {
		logger.Error(domain.CryptoAPIError, domain.UnexpectedError, "failed to parse and validate token ", err)
		return &domain.ChallengeValidationResult{
			Valid:           false,
			ValidationError: err.Error(),
		}, nil
	}

	compressedPubKey := (token.Header["kid"]).(string)
	// get challenge from repo using pub key and nonce
	challenges, err := cs.repo.ChallengeRepo.GetChallenges(compressedPubKey, claims.Id)
	if err != nil {
		logger.Error(domain.CryptoAPIError, domain.UnexpectedError, "failed to get challenge from repo; nonce: ", claims.Id)
		return &domain.ChallengeValidationResult{
			Valid: false,
		}, err
	}

	// if no challenge found in repo for the pubkey+nonce combination, it means token nonce is invalid
	if len(challenges) == 0 {
		return &domain.ChallengeValidationResult{
			Valid:           false,
			ValidationError: "invalid nonce",
		}, nil
	}

	if challenges[0].ExpiresAt < time.Now().Unix() {
		return &domain.ChallengeValidationResult{
			Valid:           false,
			ValidationError: "expired nonce",
		}, nil
	}

	return &domain.ChallengeValidationResult{
		Valid: token.Valid,
	}, nil
}

func getPublicKey(token *jwt.Token) (interface{}, error) {
	decompressedPubKey, err := decompressPublicKey((token.Header["kid"]).(string))
	if err != nil {
		logger.Error(domain.CryptoAPIError, domain.UnexpectedError, "failed to decompress public key ", err)
		return nil, err
	}
	stringKey, err := hex.DecodeString(string(decompressedPubKey))
	if err != nil {
		logger.Error(domain.CryptoAPIError, domain.UnexpectedError, "failed to decode hex public key string ", err)
		return nil, err
	}

	decodedKey, _ := pem.Decode(stringKey)
	pubKey, err := x509.ParsePKIXPublicKey(decodedKey.Bytes)
	if err != nil {
		logger.Error(domain.CryptoAPIError, domain.UnexpectedError, "failed to parse public key ", err)
		return nil, err
	}

	return pubKey, nil
}

func decompressPublicKey(compressed string) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(compressed)
	if err != nil {
		logger.Error(domain.CryptoAPIError, domain.UnexpectedError, "failed to decode compressed pub key string ", err)
		return nil, err
	}
	gzipReader, err := gzip.NewReader(bytes.NewReader(decoded))
	if err != nil {
		logger.Error(domain.CryptoAPIError, domain.UnexpectedError, "failed to create gzip reader ", err)
		return nil, err
	}
	decompressedPubKey, err := ioutil.ReadAll(gzipReader)
	if err != nil {
		logger.Error(domain.CryptoAPIError, domain.UnexpectedError, "failed to read bytes from gzip reader ", err)
		return nil, err
	}

	return decompressedPubKey, nil
}
