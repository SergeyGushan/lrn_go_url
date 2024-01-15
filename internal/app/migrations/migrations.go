package migrations

import "github.com/SergeyGushan/lrn_go_url/internal/database"

func Handle() {
	if database.Client().Ping() == nil {
		createTableUrls()
	}
}

func createTableUrls() {
	_, err := database.Client().Exec(
		"CREATE TABLE IF NOT EXISTS urls (id SERIAL PRIMARY KEY, user_id VARCHAR(36) UNIQUE NOT NULL, correlation_id  VARCHAR(255), short_url VARCHAR(255) NOT NULL, original_url VARCHAR(255) UNIQUE NOT NULL);",
	)

	if err != nil {
		panic(err)
	}
}
