package database

import (
	"database/sql"
	"github.com/SergeyGushan/lrn_go_url/cmd/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var dbClient *sql.DB

func Connect() {

	var err error
	dbClient, err = sql.Open("pgx", config.Opt.DatabaseDSN)
	if err != nil {
		panic(err)
	}
}

func Client() *sql.DB {
	return dbClient
}
