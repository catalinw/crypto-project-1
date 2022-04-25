package domain

type Challenge struct {
	PublicKey string `json:"publicKey"`
	Nonce     string `json:"nonce"`
	ExpiresAt int64  `json:"expiresAt"`
}

type ChallengeValidationResult struct {
	Valid           bool   `json:"valid"`
	ValidationError string `json:"validationError"`
}
