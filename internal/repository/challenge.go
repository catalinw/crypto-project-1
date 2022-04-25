package repository

import (
	"crypto-project-1/internal/domain"
)

type ChallengeRepo interface {
	GetChallenges(string, string) ([]*domain.Challenge, error)
	CreateChallenge(string, string, int64) (*domain.Challenge, error)
}
