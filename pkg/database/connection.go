package database

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"fmt"
)

type PgClient struct {
	pool *pgxpool.Pool
}

// NewFromEnv sets up the database connections using the configuration in the
// process's environment variables. This should be called just once per server
// instance.
func NewFromEnv(ctx context.Context, cfg *Config) (*PgClient, error) {
	pgxConfig, err := pgxpool.ParseConfig(dbDSN(cfg))
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	// BeforeAcquire is called before before a connection is acquired from the
	// pool. It must return true to allow the acquision or false to indicate that
	// the connection should be destroyed and a different connection should be
	// acquired.
	pgxConfig.BeforeAcquire = func(ctx context.Context, conn *pgx.Conn) bool {
		// Ping the connection to see if it is still valid. Ping returns an error if
		// it fails.
		return conn.Ping(ctx) == nil
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	return &PgClient{pool: pool}, nil
}

// Close releases database connections.
func (conn *PgClient) Close(ctx context.Context) {
	fmt.Println("Closing connection pool.")
	conn.pool.Close()
}

// dbDSN builds a connection string suitable for the pgx Postgres driver, using
// the values of vars.
func dbDSN(cfg *Config) string {
	vals := dbValues(cfg)
	p := make([]string, 0, len(vals))
	for k, v := range vals {
		p = append(p, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(p, " ")
}

func setIfNotEmpty(m map[string]string, key, val string) {
	if val != "" {
		m[key] = val
	}
}

func setIfPositive(m map[string]string, key string, val int) {
	if val > 0 {
		m[key] = fmt.Sprintf("%d", val)
	}
}

func dbValues(cfg *Config) map[string]string {
	p := map[string]string{}
	setIfNotEmpty(p, "dbname", cfg.Name)
	setIfNotEmpty(p, "user", cfg.User)
	setIfNotEmpty(p, "host", cfg.Host)
	setIfNotEmpty(p, "port", cfg.Port)
	setIfPositive(p, "connect_timeout", cfg.ConnectionTimeout)
	setIfNotEmpty(p, "password", cfg.Password)
	setIfNotEmpty(p, "sslmode", cfg.SSL)
	return p
}
