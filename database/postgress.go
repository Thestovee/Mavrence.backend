package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func DatabaseInit() (*pgxpool.Pool, error) {
	ctx := context.Background()

	//// Read connection string components from environment variables
	//dbHost := os.Getenv("DB_HOST")
	//dbPort := os.Getenv("DB_PORT")
	//dbUser := os.Getenv("DB_USER")
	//dbPass := os.Getenv("DB_PASSWORD")
	//dbName := os.Getenv("DB_NAME")
	//
	//if dbHost == "" || dbPort == "" || dbUser == "" || dbPass == "" || dbName == "" {
	//	return nil, fmt.Errorf("database environment variables (DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME) must be set")
	//}
	//
	//// Construct the connection string dynamically
	//connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
	//	dbUser, dbPass, dbHost, dbPort, dbName)

	connStr := fmt.Sprintf("postgres://lox:assfuck@192.168.1.23:5432//mavhealth?sslmode=disable")
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, err // Clean return of error
	}

	// Ping the database to ensure the connection is actually alive! (Good practice)
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	// Schema initialization (kept here for simplicity, but often moved to migrations)
	_, err = pool.Exec(ctx, `
            CREATE TABLE IF NOT EXISTS MavHealth (
                id      BIGSERIAL PRIMARY KEY,
                date   timestamp NOT NULL,
                mentalState int NOT NULL,
    			physicalState int NOT NULL,
    			note TEXT
            );`)
	if err != nil {
		pool.Close()
		return nil, err
	}
	_, err = pool.Exec(ctx, `
            CREATE TABLE IF NOT EXISTS WordOfTheDay (
                id      BIGSERIAL PRIMARY KEY,
                word    TEXT NOT NULL
            );`)
	if err != nil {
		pool.Close()
		return nil, err
	}
	return pool, nil
}
