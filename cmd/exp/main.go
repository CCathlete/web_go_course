package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func connectionString() string {
	PostgresConfig := struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
		SSLMode  string
	}{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSL_MODE"),
	}

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		PostgresConfig.Host, PostgresConfig.Port, PostgresConfig.User,
		PostgresConfig.Password, PostgresConfig.DBName, PostgresConfig.SSLMode,
	)
}

func main() {
	fmt.Println(connectionString())
	db, err := sql.Open("pgx", connectionString())
	if err != nil {
		panic("Error when opening db.")
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(fmt.Sprintf("Error when pinging db: %v", err))
	}
	fmt.Println("Connected!")
}
