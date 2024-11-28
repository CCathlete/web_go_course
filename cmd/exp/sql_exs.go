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

	// Inserting some data.
	type User struct {
		ID       int
		username string
	}
	type Tweet struct {
		ID, userID int
		content    string
	}
	type Like struct {
		ID, likerID, tweetID int
	}

	users := []User{
		{username: "BilboB"},
		{username: "FrodoB"},
		{username: "SamoiseG"},
		{username: "Mary"},
		{username: "Pippin"},
	}

	tweets := []Tweet{
		{userID: 1, content: "Take the ring Frodo."},
		{userID: 1, content: "I'm old."},
		{userID: 1, content: "I'm retired."},
	}

	likes := []Like{
		{likerID: 2, tweetID: 1},
		{likerID: 3, tweetID: 1},
		{likerID: 4, tweetID: 1},
		{likerID: 5, tweetID: 1},
	}

	// Inserting data into users.
	// insertDummyData := true
	insertDummyData := false
	if insertDummyData {
		for _, user := range users {
			_, err = db.Exec(`
		insert into users (username) 
		values($1);`,
				user.username)

			if err != nil {
				panic(err)
			}
		}

		fmt.Println("Users data inserted.")

		// Inserting data into tweets.
		for _, tweet := range tweets {
			_, err = db.Exec(`
		insert into tweets (userID, content) 
		values ($1, $2);`,
				tweet.userID, tweet.content)

			if err != nil {
				panic(err)
			}
		}

		fmt.Println("Tweets data inserted.")

		// Inserting data into likes.
		for _, like := range likes {
			_, err = db.Exec(`
		insert into likes (likerID, tweetID) 
		values($1, $2);`,
				like.likerID, like.tweetID)

			if err != nil {
				panic(err)
			}
		}

		fmt.Println("Likes data inserted.")

	}
	// Querying specific data.
	targetUserID := 1
	rows, err := db.Query(`
	select users.id as user_id, users.username, tweets.content
	from users
	join tweets on users.id = tweets.userID
	where users.id = $1;`,
		targetUserID)

	if err != nil {
		panic(fmt.Errorf("couldn't run query: %w", err))
	}

	for rows.Next() {
		var uid int
		var userName, currentTweet string
		err := rows.Scan(&uid, &userName, &currentTweet)
		if err != nil {
			panic(err)
		}
		fmt.Printf("User %s tweeted: %s\n", userName, currentTweet)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}
}
