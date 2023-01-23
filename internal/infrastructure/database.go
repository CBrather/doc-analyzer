package infrastructure

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/CBrather/go-auth/internal/config"
)

var db *sql.DB = nil

// GetDBClient returns the db, opening the connection first if not already done.
func GetDB() (*sql.DB, error) {
	if db != nil {
		return db, nil
	}

	db, err := sql.Open("postgres", buildConnectionString())
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func buildConnectionString() string {
	env := config.GetEnvironment()

	connectionString := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		env.Database.Host, env.Database.Port, env.Database.Name, env.Database.User, env.Database.Password, env.Database.SslMode,
	)

	return connectionString
}
