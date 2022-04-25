package public

type CreateChallengeRequestBody struct {
	PubKey string `json:"pubKey"`
}

type VerifyChallengeRequestBody struct {
	Token string `json:"token"`
}
