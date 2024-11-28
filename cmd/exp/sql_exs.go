package main

import (
	"database/sql"
	"fmt"
	"webGo/models"
)

func sqlExs() {
	db := Must(models.ConnectToDB()).(*sql.DB)
	defer db.Close()

	// Creating a table.
	Must(db.Exec(`
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
	`))

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
			Must(db.Exec(`
		insert into users (username) 
		values($1);`,
				user.username))
		}

		fmt.Println("Users data inserted.")

		// Inserting data into tweets.
		for _, tweet := range tweets {
			Must(db.Exec(`
		insert into tweets (userID, content) 
		values ($1, $2);`,
				tweet.userID, tweet.content))
		}

		fmt.Println("Tweets data inserted.")

		// Inserting data into likes.
		for _, like := range likes {
			Must(db.Exec(`
		insert into likes (likerID, tweetID) 
		values($1, $2);`,
				like.likerID, like.tweetID))
		}

		fmt.Println("Likes data inserted.")

	}
	// Querying specific data.
	targetUserID := 1
	rows := Must(db.Query(`
	select users.id as user_id, users.username, tweets.content
	from users
	join tweets on users.id = tweets.userID
	where users.id = $1;`,
		targetUserID)).(*sql.Rows)

	for rows.Next() {
		var uid int
		var userName, currentTweet string
		Must(nil, rows.Scan(&uid, &userName, &currentTweet))
		fmt.Printf("User %s tweeted: %s\n", userName, currentTweet)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}
}
