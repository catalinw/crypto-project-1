package tests

import (
	"bytes"
	"crypto-project-1/public"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/cucumber/godog"
	"github.com/gofrs/uuid"
	"github.com/joho/godotenv"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
	// use pq as a library to create postgres client
	_ "github.com/lib/pq"
)

const (
	dbPortVar          = "PGPORT"
	dbHostVar          = "PGHOST"
	dbNameVar          = "PGDATABASE"
	dbUserVar          = "PGUSER"
	dbPassVar          = "PGPASSWORD"
	challengeTableName = "challenge"
)

type challengeTest struct {
	db                 *sql.DB
	publicKey          string
	nonce              string
	expiresAt          int64
	token              string
	verifyResponseBody []byte
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	// read env variables from .env file
	_ = godotenv.Overload()

	database, err := newDB()
	if err != nil {
		panic(fmt.Sprintf("TEST FAILED: failed to create db client, err: %s", err))
	}
	err = database.Ping()
	if err != nil {
		panic(fmt.Sprintf("TEST FAILED: failed to ping db, err: %s", err))
	}

	challengeTest := challengeTest{
		publicKey: "H4sIAAAAAAAA/4SQQU4EMQwEv5TY1e34OZmdyf+fgBaEQFyQb6U6uCvu7yMQRfPE0JAIXjQZgwupf8xxj82NfVWhKjWLy0ebru2dg2Rqquk/nDd3gqaLiZBcGccnFVdJ0yG8eWU6var98EqqtFxfsiL7/UPNbEnKRad8I5+yZOdKa2pn+eIQtDe7qqaf2pDTx7eny0IsVRXLEJw0EZfBng7fRmRrcGq58s7P/b+6iQf+axbjAwAA//8BAAD//0A4Ig9qAQAA",
		nonce:     "4b8b3887-e113-4e27-adb4-06f9aa66c395",
		expiresAt: 1893492000,
		token:     "eyJhbGciOiJFUzI1NiIsImtpZCI6Ikg0c0lBQUFBQUFBQS80U1FRVTRFTVF3RXY1VFkxZTM0T1ptZHlmK2ZnQmFFUUZ5UWI2VTZ1Q3Z1N3lNUVJmUEUwSkFJWGpRWmd3dXBmOHh4ajgyTmZWV2hLaldMeTBlYnJ1MmRnMlJxcXVrL25EZDNncWFMaVpCY0djY25GVmRKMHlHOGVXVTZ2YXI5OEVxcXRGeGZzaUw3L1VQTmJFbktSYWQ4STUreVpPZEthMnBuK2VJUXREZTdxcWFmMnBEVHg3ZW55MElzVlJYTEVKdzBFWmZCbmc3ZlJtUnJjR3E1OHM3UC9iKzZpUWYrYXhiakF3QUEvLzhCQUFELy8wQTRJZzlxQVFBQSIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ3aGVsdGVlIiwiZXhwIjoxODkzNDkyMDAwLCJqdGkiOiI0YjhiMzg4Ny1lMTEzLTRlMjctYWRiNC0wNmY5YWE2NmMzOTUiLCJpYXQiOjE2NTA5MTc0MzQsIm5iZiI6MTY1MDkxNzQzNH0.qjSsBaeDtb4Iesraq6McH-M9Iqh7zZSP_bSuJVg9dvbSKjo_WwQQLpJIi20S_vrrUCiJI3WYDq31SicdnXLbxg",
		db:        database,
	}

	// these tests will actually delete data from the database, please run the tests only on testing envs
	ctx.Step(`^a clean database$`, challengeTest.aCleanDatabase)
	ctx.Step(`^I send a request to create a challenge$`, challengeTest.iSendARequestToCreateAChallenge)
	ctx.Step(`^I wait for the request to be processed$`, challengeTest.iWaitForTheRequestToBeProcessed)
	ctx.Step(`^the challenge should be created and valid$`, challengeTest.theChallengeShouldBeCreatedAndValid)

	ctx.Step(`^a challenge that was previously created$`, challengeTest.aChallengeThatWasPreviouslyCreated)
	ctx.Step(`^I send a request to validate a challenge$`, challengeTest.iSendARequestToValidateAChallenge)
	ctx.Step(`^the challenge should be validated successfully$`, challengeTest.theChallengeShouldBeValidatedSuccessfully)

}

func (ct *challengeTest) aCleanDatabase() error {
	qb := ct.dbQueryBuilder().
		Delete(challengeTableName).
		Where(squirrel.Eq{"public_key": ct.publicKey})

	_, err := qb.Exec()
	if err != nil {
		return fmt.Errorf("TEST FAILED: failed to clean database, err: %w", err)
	}

	return nil
}

func (ct *challengeTest) iSendARequestToCreateAChallenge() error {
	createChallengeBody := &public.CreateChallengeRequestBody{
		PubKey: ct.publicKey,
	}
	jsonBody, err := json.Marshal(createChallengeBody)
	if err != nil {
		return fmt.Errorf("TEST FAILED: failed to marshal create challenge request body, err: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, "http://localhost:7777/v1/challenge", bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("TEST FAILED: failed to create http request, err: %w", err)
	}
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("TEST FAILED: failed to send request to API endpoint, err: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("TEST FAILED: expected 200 OK from endpoint, received http status code: %d", response.StatusCode)
	}

	return nil
}

func (ct *challengeTest) iWaitForTheRequestToBeProcessed() error {
	time.Sleep(time.Second * 1)

	return nil
}

func (ct *challengeTest) theChallengeShouldBeCreatedAndValid() error {
	qb := ct.dbQueryBuilder().
		Select("public_key", "nonce", "expires_at").
		From(challengeTableName).
		Where(squirrel.Eq{"public_key": ct.publicKey})
	rows, err := qb.Query()
	if err != nil {
		return fmt.Errorf("TEST FAILED: failed to create db query, err: %w", err)
	}

	type challenge struct {
		publicKey string
		nonce     string
		expiresAt int64
	}

	var dbResults []*challenge
	for rows.Next() {
		var c challenge
		err = rows.Scan(&c.publicKey, &c.nonce, &c.expiresAt)
		if err != nil {
			return fmt.Errorf("TEST FAILED: failed to execute db query, err: %w", err)
		}

		dbResults = append(dbResults, &c)
	}

	if len(dbResults) == 0 {
		return fmt.Errorf("TEST FAILED: challenge was not created in db")
	}

	if dbResults[0].expiresAt < time.Now().Unix() {
		return fmt.Errorf("TEST FAILED: created challenge has expired nonce")
	}

	if _, err := uuid.FromString(dbResults[0].nonce); err != nil {
		return fmt.Errorf("TEST FAILED: created challenge contains nonce that is not a valid uuid")
	}

	return nil
}

func (ct *challengeTest) aChallengeThatWasPreviouslyCreated() error {
	qb := ct.dbQueryBuilder().
		Insert(challengeTableName).
		Columns("public_key", "nonce", "expires_at").
		Values(ct.publicKey, ct.nonce, ct.expiresAt).
		Suffix("RETURNING nonce")

	var createdNonce string
	err := qb.QueryRow().Scan(&createdNonce)
	if err != nil {
		return fmt.Errorf("TEST FAILED: failed to insert challenge into db")
	}

	return nil
}

func (ct *challengeTest) iSendARequestToValidateAChallenge() error {
	verifyChallengeBody := &public.VerifyChallengeRequestBody{
		Token: ct.token,
	}
	jsonBody, err := json.Marshal(verifyChallengeBody)
	if err != nil {
		return fmt.Errorf("TEST FAILED: failed to marshal verify challenge request body, err: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, "http://localhost:7777/v1/verify-challenge", bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("TEST FAILED: failed to create http request, err: %w", err)
	}
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("TEST FAILED: failed to send request to API endpoint, err: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("TEST FAILED: expected 200 OK from endpoint, received http status code: %d", response.StatusCode)
	}
	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("TEST FAILED: failed to parse response body: %w", err)
	}

	ct.verifyResponseBody = respBody

	return nil
}

func (ct *challengeTest) theChallengeShouldBeValidatedSuccessfully() error {
	type apiResponse struct {
		Result struct {
			Valid           bool   `json:"valid"`
			ValidationError string `json:"validationError"`
		} `json:"result"`
	}

	response := &apiResponse{}
	if err := json.Unmarshal(ct.verifyResponseBody, response); err != nil {
		return fmt.Errorf("TEST FAILED: failed to unmarshal response body: %w", err)
	}

	if !response.Result.Valid {
		return fmt.Errorf("TEST FAILED: token validation failed, validation error: %s", response.Result.ValidationError)
	}

	return nil
}

func (ct *challengeTest) dbQueryBuilder() squirrel.StatementBuilderType {
	return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).RunWith(ct.db)
}

func newDB() (*sql.DB, error) {
	host, found := os.LookupEnv(dbHostVar)
	if !found {
		err := fmt.Errorf("missing env variable %s", dbHostVar)
		return nil, err
	}
	p, found := os.LookupEnv(dbPortVar)
	if !found {
		err := fmt.Errorf("missing env variable %s", dbPortVar)
		return nil, err
	}
	port, err := strconv.Atoi(p)
	if err != nil {
		return nil, err
	}
	dbname, found := os.LookupEnv(dbNameVar)
	if !found {
		err := fmt.Errorf("missing env variable %s", dbNameVar)
		return nil, err
	}
	user, found := os.LookupEnv(dbUserVar)
	if !found {
		err := fmt.Errorf("missing env variable %s", dbUserVar)
		return nil, err
	}
	password, found := os.LookupEnv(dbPassVar)
	if !found {
		err := fmt.Errorf("missing env variable %s", dbPassVar)
		return nil, err
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		host, port, user, dbname, password)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	return db, nil
}
