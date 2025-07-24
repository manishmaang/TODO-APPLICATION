package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func init() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	// Parse connection string into config
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatal("Failed to parse DATABASE_URL:", err)
	}

	// Customize connection pool settings
	config.MaxConns = 10                       // max number of connections in the pool
	config.MinConns = 2                        // minimum number of connections always kept alive
	config.MaxConnIdleTime = 5 * time.Minute   // close idle connections after 5 minutes
	config.MaxConnLifetime = 30 * time.Minute  // close connections older than 30 minutes
	config.HealthCheckPeriod = 30 * time.Second // periodically ping idle connections

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatal("Unable to connect to the DB:", err)
	}

	if err = db.Ping(ctx); err != nil {
		log.Fatal("Unable to ping the DB:", err)
	}

	DB = db
	fmt.Println("Successfully connected to the Database with custom pool settings")
}
