package repository

import (
	"crypto-project-1/internal/domain"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	// use pq as a library to create postgres client
	_ "github.com/lib/pq"
	"os"
	"strconv"
)

const (
	dbPortVar = "PGPORT"
	dbHostVar = "PGHOST"
	dbNameVar = "PGDATABASE"
	dbUserVar = "PGUSER"
	dbPassVar = "PGPASSWORD"
)

var db *sql.DB

func dbQueryBuilder() squirrel.StatementBuilderType {
	return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).RunWith(db)
}

func NewDB() (*sql.DB, error) {
	host, found := os.LookupEnv(dbHostVar)
	if !found {
		err := fmt.Errorf("%s missing env variable %s", domain.CryptoAPIError, dbHostVar)
		return nil, err
	}
	p, found := os.LookupEnv(dbPortVar)
	if !found {
		err := fmt.Errorf("%s missing env variable %s", domain.CryptoAPIError, dbPortVar)
		return nil, err
	}
	port, err := strconv.Atoi(p)
	if err != nil {
		return nil, err
	}
	dbname, found := os.LookupEnv(dbNameVar)
	if !found {
		err := fmt.Errorf("%s missing env variable %s", domain.CryptoAPIError, dbNameVar)
		return nil, err
	}
	user, found := os.LookupEnv(dbUserVar)
	if !found {
		err := fmt.Errorf("%s missing env variable %s", domain.CryptoAPIError, dbUserVar)
		return nil, err
	}
	password, found := os.LookupEnv(dbPassVar)
	if !found {
		err := fmt.Errorf("%s missing env variable %s", domain.CryptoAPIError, dbPassVar)
		return nil, err
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		host, port, user, dbname, password)
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	return db, nil
}
