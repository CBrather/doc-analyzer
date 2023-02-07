package infrastructure

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"

	"github.com/CBrather/analyzer/internal/config"
)

var db *sql.DB = nil

// GetDBClient returns the db, opening the connection first if not already done.
func GetDB() (*sql.DB, error) {
	if db != nil {
		return db, nil
	}

	env := config.GetEnvironment()

	db, err := otelsql.Open("postgres", buildConnectionString(env), otelsql.WithAttributes(
		semconv.DBSystemPostgreSQL,
		semconv.DBNameKey.String(env.Database.Name),
	))
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	otelsql.ReportDBStatsMetrics(db, otelsql.WithAttributes(semconv.DBSystemPostgreSQL))

	return db, nil
}

func buildConnectionString(env *config.EnvConfig) string {
	connectionString := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		env.Database.Host, env.Database.Port, env.Database.Name, env.Database.User, env.Database.Password, env.Database.SslMode,
	)

	return connectionString
}
