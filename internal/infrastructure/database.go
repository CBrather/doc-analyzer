package infrastructure

import (
	"database/sql"
	"fmt"

	"github.com/XSAM/otelsql"
	_ "github.com/lib/pq"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"

	"github.com/CBrather/go-auth/internal/config"
)

var db *sql.DB = nil

// GetDBClient returns the db, opening the connection first if not already done.
func GetDB() (*sql.DB, error) {
	if db != nil {
		return db, nil
	}

	db, err := otelsql.Open("postgres", buildConnectionString(), otelsql.WithAttributes(
		semconv.DBSystemPostgreSQL,
	))
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	if err = otelsql.RegisterDBStatsMetrics(db, otelsql.WithAttributes(semconv.DBSystemPostgreSQL)); err != nil {
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
