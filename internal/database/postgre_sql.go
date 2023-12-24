package database

import (
	"database/sql"
	"github.com/SergeyGushan/lrn_go_url/cmd/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var DBClient *sql.DB

func Connect() {

	var err error
	DBClient, err = sql.Open("pgx", config.Opt.DatabaseDSN)
	if err != nil {
		panic(err)
	}
}
