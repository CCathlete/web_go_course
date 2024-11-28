package main

import (
	"database/sql"
	"fmt"
	"os"
	"webGo/models"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	db := connectToDB()
	defer db.Close()
	userService := models.UserService{
		DB: db,
	}
	Must(userService.Create("kcat@example.com", "secret pass"))
}

func connectToDB() *sql.DB {
	db, err := sql.Open("pgx", connectionString())
	if err != nil {
		panic("Error when opening db.")
	}
	Must(nil, db.Ping())

	return db
}

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

func Must(value any, err error) any {
	if err != nil {
		panic(err)
	}

	return value
}
