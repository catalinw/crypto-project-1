package repository

import (
	"crypto-project-1/internal/domain"
	"github.com/Masterminds/squirrel"
	logger "github.com/sirupsen/logrus"
)

const (
	challengeTableName = "challenge"
)

type challengeDbRepo struct{}

func (db *challengeDbRepo) GetChallenges(pubKey, nonce string) ([]*domain.Challenge, error) {
	qb := dbQueryBuilder().
		Select("public_key", "nonce", "expires_at").
		From(challengeTableName).
		Where(squirrel.And{
			squirrel.Eq{"public_key": pubKey},
			squirrel.Eq{"nonce": nonce},
		})
	rows, err := qb.Query()
	if err != nil {
		logger.Error(domain.CryptoAPIError, domain.UnexpectedError, "failed to create get query ", err)
		return nil, err
	}

	var challenges []*domain.Challenge
	for rows.Next() {
		var challenge domain.Challenge
		err = rows.Scan(&challenge.PublicKey, &challenge.Nonce, &challenge.ExpiresAt)
		if err != nil {
			logger.Error(domain.CryptoAPIError, domain.UnexpectedError, "failed to execute get query ", err)
			return nil, err
		}

		challenges = append(challenges, &challenge)
	}

	return challenges, nil
}

func (db *challengeDbRepo) CreateChallenge(pubKey, nonce string, expiresAt int64) (*domain.Challenge, error) {
	qb := dbQueryBuilder().
		Insert(challengeTableName).
		Columns("public_key", "nonce", "expires_at").
		Values(pubKey, nonce, expiresAt).
		Suffix("RETURNING nonce")

	var createdNonce string
	err := qb.QueryRow().Scan(&createdNonce)
	if err != nil {
		logger.Error(domain.CryptoAPIError, domain.UnexpectedError, "failed to execute insert query ", err)
		return nil, err
	}

	return &domain.Challenge{
		PublicKey: pubKey,
		Nonce:     createdNonce,
		ExpiresAt: expiresAt,
	}, nil
}
