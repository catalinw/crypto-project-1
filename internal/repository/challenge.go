package repository

import (
	"crypto-project-1/internal/domain"
)

//go:generate mockgen -package=mock_repository -destination=./mock_repository/challenge.go -source=challenge.go
type ChallengeRepository interface {
	GetChallenges(string, string) ([]*domain.Challenge, error)
	CreateChallenge(string, string, int64) (*domain.Challenge, error)
}
