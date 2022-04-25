package repository

type Repo struct {
	ChallengeRepo ChallengeRepo
}

func NewRepo() *Repo {
	return &Repo{
		ChallengeRepo: &challengeDbRepo{},
	}
}
