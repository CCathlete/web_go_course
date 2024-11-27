package main

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func sqlExs() {
	db, err := sql.Open("pgx", connectionString())
	if err != nil {
		panic("Error when opening db.")
	}
	defer db.Close()

	// Creating a table.
	_, err = db.Exec(`
		create table if not exists users (
			id serial primary key,
			username text unique not null
		);

		create table if not exists tweets (
			id serial primary key,
			userID int not null,
			content text,
			foreign key (userID) references users(id)
		);

		create table if not exists likes (
			id serial primary key,
			tweetID int not null,
			likerID int,
			foreign key (likerID) references users(id),
			foreign key (tweetID) references tweets(id)
		);
	`)
	if err != nil {
		panic(err)
	}

	fmt.Println("Tables created.")

	// Insert some data.
	type user struct {
		ID       int
		username string
	}
	type tweet struct {
		ID, userID int
		content    string
	}
	type like struct {
		ID, likerID, tweetID int
	}

	usersDetails := []user{
		{username: "BilboB"},
		{username: "FrodoB"},
		{username: "SamoiseG"},
		{username: "Mary"},
		{username: "Pippin"},
	}

	userTweets := []tweet{
		{userID: 1, content: "Take the ring Frodo."},
		{userID: 1, content: "I'm old."},
		{userID: 1, content: "I'm retired."},
	}

	tweetLikes := []like{
		{likerID: 2, tweetID: 1},
		{likerID: 3, tweetID: 1},
		{likerID: 4, tweetID: 1},
		{likerID: 5, tweetID: 1},
	}

}

func main() {
	sqlExs()
}
