package migrations

import "github.com/SergeyGushan/lrn_go_url/internal/database"

func Handle() {
	if database.DBClient.Ping() == nil {
		createTableUrls()
	}
}

func createTableUrls() {
	_, err := database.DBClient.Exec(
		"CREATE TABLE IF NOT EXISTS urls (id SERIAL PRIMARY KEY, short_url VARCHAR(255) NOT NULL, original_url VARCHAR(255) NOT NULL);",
	)

	if err != nil {
		panic(err)
	}
}
