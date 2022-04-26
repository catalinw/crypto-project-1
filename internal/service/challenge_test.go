package service_test

import (
	"crypto-project-1/internal/domain"
	"crypto-project-1/internal/repository"
	"crypto-project-1/internal/repository/mock_repository"
	"crypto-project-1/internal/service"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestChallengeService_CreateChallenge(t *testing.T) {
	type args struct {
		publicKey string
		now       func() time.Time
	}

	type expected struct {
		repoCreateIsCalled bool
		expiresAt          int64
		publicKey          string
		nonce              string
		errorIsReturned    bool
	}

	timeNow := time.Now()

	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "create challenge successfully using valid public key",
			args: args{
				publicKey: "H4sIAAAAAAAA/4SQQU4EMQwEv5TY1e34OZmdyf+fgBaEQFyQb6U6uCvu7yMQRfPE0JAIXjQZgwupf8xxj82NfVWhKjWLy0ebru2dg2Rqquk/nDd3gqaLiZBcGccnFVdJ0yG8eWU6var98EqqtFxfsiL7/UPNbEnKRad8I5+yZOdKa2pn+eIQtDe7qqaf2pDTx7eny0IsVRXLEJw0EZfBng7fRmRrcGq58s7P/b+6iQf+axbjAwAA//8BAAD//0A4Ig9qAQAA",
				now: func() time.Time {
					return timeNow
				},
			},
			expected: expected{
				repoCreateIsCalled: true,
				expiresAt:          timeNow.Add(time.Minute * 5).Unix(),
				publicKey:          "H4sIAAAAAAAA/4SQQU4EMQwEv5TY1e34OZmdyf+fgBaEQFyQb6U6uCvu7yMQRfPE0JAIXjQZgwupf8xxj82NfVWhKjWLy0ebru2dg2Rqquk/nDd3gqaLiZBcGccnFVdJ0yG8eWU6var98EqqtFxfsiL7/UPNbEnKRad8I5+yZOdKa2pn+eIQtDe7qqaf2pDTx7eny0IsVRXLEJw0EZfBng7fRmRrcGq58s7P/b+6iQf+axbjAwAA//8BAAD//0A4Ig9qAQAA",
				nonce:              uuid.NewString(),
				errorIsReturned:    false,
			},
		},
		{
			name: "create challenge fails using invalid public key",
			args: args{
				publicKey: "testingCataPublicKey",
				now: func() time.Time {
					return timeNow
				},
			},
			expected: expected{
				repoCreateIsCalled: false,
				errorIsReturned:    true,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepo := mock_repository.NewMockChallengeRepository(ctrl)

			if test.expected.repoCreateIsCalled {
				mockRepo.EXPECT().CreateChallenge(test.expected.publicKey, gomock.Any(), test.expected.expiresAt).
					Return(&domain.Challenge{
						PublicKey: test.expected.publicKey,
						ExpiresAt: test.expected.expiresAt,
						Nonce:     test.expected.nonce,
					}, nil)
			}

			repo := repository.NewRepository(mockRepo)
			challengeService := service.NewChallengeService(repo, test.args.now)
			challenge, err := challengeService.CreateChallenge(test.args.publicKey)
			if test.expected.errorIsReturned {
				assert.Error(t, err)

				return
			}

			assert.NoError(t, err)
			// check nonce is a valid uuid
			_, err = uuid.Parse(challenge.Nonce)
			assert.NoError(t, err)
			assert.Equal(t, test.expected.publicKey, challenge.PublicKey)
			assert.Equal(t, test.expected.expiresAt, challenge.ExpiresAt)
		})
	}
}

func TestChallengeService_VerifyChallenge(t *testing.T) {
	type args struct {
		signedToken            string
		now                    func() time.Time
		repoReturnedChallenges []*domain.Challenge
	}

	type expected struct {
		repoGetIsCalled bool
		publicKey       string
		nonce           string
		tokenIsValid    bool
		errorIsReturned bool
		validationError string
	}

	timeNow := time.Now()

	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "verify challenge successfully using valid signed token",
			args: args{
				signedToken: "eyJhbGciOiJFUzI1NiIsImtpZCI6Ikg0c0lBQUFBQUFBQS80U1FRVTRFTVF3RXY1VFkxZTM0T1ptZHlmK2ZnQmFFUUZ5UWI2VTZ1Q3Z1N3lNUVJmUEUwSkFJWGpRWmd3dXBmOHh4ajgyTmZWV2hLaldMeTBlYnJ1MmRnMlJxcXVrL25EZDNncWFMaVpCY0djY25GVmRKMHlHOGVXVTZ2YXI5OEVxcXRGeGZzaUw3L1VQTmJFbktSYWQ4STUreVpPZEthMnBuK2VJUXREZTdxcWFmMnBEVHg3ZW55MElzVlJYTEVKdzBFWmZCbmc3ZlJtUnJjR3E1OHM3UC9iKzZpUWYrYXhiakF3QUEvLzhCQUFELy8wQTRJZzlxQVFBQSIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ3aGVsdGVlIiwiZXhwIjoxODkzNDkyMDAwLCJqdGkiOiI0YjhiMzg4Ny1lMTEzLTRlMjctYWRiNC0wNmY5YWE2NmMzOTUiLCJpYXQiOjE2NTA5MTc0MzQsIm5iZiI6MTY1MDkxNzQzNH0.qjSsBaeDtb4Iesraq6McH-M9Iqh7zZSP_bSuJVg9dvbSKjo_WwQQLpJIi20S_vrrUCiJI3WYDq31SicdnXLbxg",
				now: func() time.Time {
					return timeNow
				},
				repoReturnedChallenges: []*domain.Challenge{
					{
						PublicKey: "H4sIAAAAAAAA/4SQQU4EMQwEv5TY1e34OZmdyf+fgBaEQFyQb6U6uCvu7yMQRfPE0JAIXjQZgwupf8xxj82NfVWhKjWLy0ebru2dg2Rqquk/nDd3gqaLiZBcGccnFVdJ0yG8eWU6var98EqqtFxfsiL7/UPNbEnKRad8I5+yZOdKa2pn+eIQtDe7qqaf2pDTx7eny0IsVRXLEJw0EZfBng7fRmRrcGq58s7P/b+6iQf+axbjAwAA//8BAAD//0A4Ig9qAQAA",
						Nonce:     "4b8b3887-e113-4e27-adb4-06f9aa66c395",
						ExpiresAt: timeNow.Add(time.Minute * 5).Unix(),
					},
				},
			},
			expected: expected{
				repoGetIsCalled: true,
				publicKey:       "H4sIAAAAAAAA/4SQQU4EMQwEv5TY1e34OZmdyf+fgBaEQFyQb6U6uCvu7yMQRfPE0JAIXjQZgwupf8xxj82NfVWhKjWLy0ebru2dg2Rqquk/nDd3gqaLiZBcGccnFVdJ0yG8eWU6var98EqqtFxfsiL7/UPNbEnKRad8I5+yZOdKa2pn+eIQtDe7qqaf2pDTx7eny0IsVRXLEJw0EZfBng7fRmRrcGq58s7P/b+6iQf+axbjAwAA//8BAAD//0A4Ig9qAQAA",
				nonce:           "4b8b3887-e113-4e27-adb4-06f9aa66c395",
				errorIsReturned: false,
				tokenIsValid:    true,
				validationError: "",
			},
		},
		{
			name: "verify challenge fails using nonce that is not stored in repo",
			args: args{
				signedToken: "eyJhbGciOiJFUzI1NiIsImtpZCI6Ikg0c0lBQUFBQUFBQS80U1FRVTRFTVF3RXY1VFkxZTM0T1ptZHlmK2ZnQmFFUUZ5UWI2VTZ1Q3Z1N3lNUVJmUEUwSkFJWGpRWmd3dXBmOHh4ajgyTmZWV2hLaldMeTBlYnJ1MmRnMlJxcXVrL25EZDNncWFMaVpCY0djY25GVmRKMHlHOGVXVTZ2YXI5OEVxcXRGeGZzaUw3L1VQTmJFbktSYWQ4STUreVpPZEthMnBuK2VJUXREZTdxcWFmMnBEVHg3ZW55MElzVlJYTEVKdzBFWmZCbmc3ZlJtUnJjR3E1OHM3UC9iKzZpUWYrYXhiakF3QUEvLzhCQUFELy8wQTRJZzlxQVFBQSIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ3aGVsdGVlIiwiZXhwIjoxODkzNDkyMDAwLCJqdGkiOiI0YjhiMzg4Ny1lMTEzLTRlMjctYWRiNC0wNmY5YWE2NmMzOTUiLCJpYXQiOjE2NTA5MTc0MzQsIm5iZiI6MTY1MDkxNzQzNH0.qjSsBaeDtb4Iesraq6McH-M9Iqh7zZSP_bSuJVg9dvbSKjo_WwQQLpJIi20S_vrrUCiJI3WYDq31SicdnXLbxg",
				now: func() time.Time {
					return timeNow
				},
				repoReturnedChallenges: []*domain.Challenge{},
			},
			expected: expected{
				repoGetIsCalled: true,
				publicKey:       "H4sIAAAAAAAA/4SQQU4EMQwEv5TY1e34OZmdyf+fgBaEQFyQb6U6uCvu7yMQRfPE0JAIXjQZgwupf8xxj82NfVWhKjWLy0ebru2dg2Rqquk/nDd3gqaLiZBcGccnFVdJ0yG8eWU6var98EqqtFxfsiL7/UPNbEnKRad8I5+yZOdKa2pn+eIQtDe7qqaf2pDTx7eny0IsVRXLEJw0EZfBng7fRmRrcGq58s7P/b+6iQf+axbjAwAA//8BAAD//0A4Ig9qAQAA",
				nonce:           "4b8b3887-e113-4e27-adb4-06f9aa66c395",
				validationError: "invalid nonce",
				errorIsReturned: false,
				tokenIsValid:    false,
			},
		},
		{
			name: "verify challenge fails using expired nonce",
			args: args{
				signedToken: "eyJhbGciOiJFUzI1NiIsImtpZCI6Ikg0c0lBQUFBQUFBQS80U1FRVTRFTVF3RXY1VFkxZTM0T1ptZHlmK2ZnQmFFUUZ5UWI2VTZ1Q3Z1N3lNUVJmUEUwSkFJWGpRWmd3dXBmOHh4ajgyTmZWV2hLaldMeTBlYnJ1MmRnMlJxcXVrL25EZDNncWFMaVpCY0djY25GVmRKMHlHOGVXVTZ2YXI5OEVxcXRGeGZzaUw3L1VQTmJFbktSYWQ4STUreVpPZEthMnBuK2VJUXREZTdxcWFmMnBEVHg3ZW55MElzVlJYTEVKdzBFWmZCbmc3ZlJtUnJjR3E1OHM3UC9iKzZpUWYrYXhiakF3QUEvLzhCQUFELy8wQTRJZzlxQVFBQSIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ3aGVsdGVlIiwiZXhwIjoxODkzNDkyMDAwLCJqdGkiOiI0YjhiMzg4Ny1lMTEzLTRlMjctYWRiNC0wNmY5YWE2NmMzOTUiLCJpYXQiOjE2NTA5MTc0MzQsIm5iZiI6MTY1MDkxNzQzNH0.qjSsBaeDtb4Iesraq6McH-M9Iqh7zZSP_bSuJVg9dvbSKjo_WwQQLpJIi20S_vrrUCiJI3WYDq31SicdnXLbxg",
				now: func() time.Time {
					return timeNow
				},
				repoReturnedChallenges: []*domain.Challenge{
					{
						PublicKey: "H4sIAAAAAAAA/4SQQU4EMQwEv5TY1e34OZmdyf+fgBaEQFyQb6U6uCvu7yMQRfPE0JAIXjQZgwupf8xxj82NfVWhKjWLy0ebru2dg2Rqquk/nDd3gqaLiZBcGccnFVdJ0yG8eWU6var98EqqtFxfsiL7/UPNbEnKRad8I5+yZOdKa2pn+eIQtDe7qqaf2pDTx7eny0IsVRXLEJw0EZfBng7fRmRrcGq58s7P/b+6iQf+axbjAwAA//8BAAD//0A4Ig9qAQAA",
						Nonce:     "4b8b3887-e113-4e27-adb4-06f9aa66c395",
						ExpiresAt: timeNow.Add(time.Minute * -5).Unix(),
					},
				},
			},
			expected: expected{
				repoGetIsCalled: true,
				publicKey:       "H4sIAAAAAAAA/4SQQU4EMQwEv5TY1e34OZmdyf+fgBaEQFyQb6U6uCvu7yMQRfPE0JAIXjQZgwupf8xxj82NfVWhKjWLy0ebru2dg2Rqquk/nDd3gqaLiZBcGccnFVdJ0yG8eWU6var98EqqtFxfsiL7/UPNbEnKRad8I5+yZOdKa2pn+eIQtDe7qqaf2pDTx7eny0IsVRXLEJw0EZfBng7fRmRrcGq58s7P/b+6iQf+axbjAwAA//8BAAD//0A4Ig9qAQAA",
				nonce:           "4b8b3887-e113-4e27-adb4-06f9aa66c395",
				validationError: "expired nonce",
				errorIsReturned: false,
				tokenIsValid:    false,
			},
		},
		{
			name: "verify challenge fails using token without public key",
			args: args{
				signedToken: "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ3aGVsdGVlIiwiZXhwIjoxODkzNDkyMDAwLCJqdGkiOiI0YjhiMzg4Ny1lMTEzLTRlMjctYWRiNC0wNmY5YWE2NmMzOTUiLCJpYXQiOjE2NTA5MTc0MzQsIm5iZiI6MTY1MDkxNzQzNH0.unga02D03TwdYQWxwMUbVCRFOGnqQxd7RPUpessP1O8amTcbsRDVrm4aqxVtiZQb6MzYIUjnc8Z8GWpiNApHWw",
				now: func() time.Time {
					return timeNow
				},
				repoReturnedChallenges: []*domain.Challenge{},
			},
			expected: expected{
				repoGetIsCalled: false,
				validationError: "public key header not found",
				errorIsReturned: false,
				tokenIsValid:    false,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepo := mock_repository.NewMockChallengeRepository(ctrl)

			if test.expected.repoGetIsCalled {
				mockRepo.EXPECT().GetChallenges(test.expected.publicKey, test.expected.nonce).Return(test.args.repoReturnedChallenges, nil)
			}

			repo := repository.NewRepository(mockRepo)
			challengeService := service.NewChallengeService(repo, test.args.now)
			validationResult, err := challengeService.VerifyChallenge(test.args.signedToken)
			if test.expected.errorIsReturned {
				assert.Error(t, err)

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, test.expected.tokenIsValid, validationResult.Valid)
			assert.Equal(t, test.expected.validationError, validationResult.ValidationError)
		})
	}
}
