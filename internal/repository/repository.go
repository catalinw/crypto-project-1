package repository

type Repository struct {
	ChallengeRepo ChallengeRepository
}

func NewRepository(challengeRepository ChallengeRepository) *Repository {
	return &Repository{
		ChallengeRepo: challengeRepository,
	}
}
