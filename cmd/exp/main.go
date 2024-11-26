package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
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

	// Creating a table.
	_, err = db.Exec(`
		create table if not exists users (
			id serial primary key,
			name text,
			email text unique not null
		);

		create table if not exists orders (
			id serial primary key,
			user_id int not null,
			amount int,
			description text
		);
	`)
	if err != nil {
		panic(err)
	}

	fmt.Println("Tables created.")

	// Insert some data.
	usersDetails := []struct{ name, email string }{
		{name: "Ken Cat", email: "ccat@example.com"},
		{name: "Jon Calhoun", email: "JC@example.com"},
		{name: "Missy Elliott", email: "doitdoitdoitdoit@example.com"},
		{name: "Aero Smith", email: "dreamon@example.com"},
	}

	for _, user := range usersDetails {
		_, err = db.Exec(`
		insert into users (name, email) values
		($1, $2);
		`, user.name, user.email)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Values inserted.")

	_, err = db.Exec(`
		select * from users;
		-- drop table users;
	`)
	if err != nil {
		panic(err)
	}
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
