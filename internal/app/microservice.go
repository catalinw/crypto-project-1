package app

import "crypto-project-1/internal/service"

type CryptoMicroservice struct {
	challengeService service.ChallengeService
}

func NewCryptoMicroservice(challengeService service.ChallengeService) *CryptoMicroservice {
	return &CryptoMicroservice{
		challengeService: challengeService,
	}
}
