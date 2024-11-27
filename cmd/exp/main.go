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

	// id := 4
	// var name, email string
	// row := db.QueryRow(`
	// select name, email
	// from users
	// where id=$1;
	// `, id)
	//
	// switch err := row.Scan(&name, &email); err {
	// case sql.ErrNoRows:
	// fmt.Printf("No matches were found: %v", err)
	// case nil:
	// break
	// default:
	// panic(err)
	// }
	//
	// fmt.Println(name, email)

	// Creating orders for later.
	// userID := 4
	// for i := 1; i <= 5; i++ {
	// 	amount := i * 100
	// 	desc := fmt.Sprintf("Fake order %d", amount)
	// 	_, err := db.Exec(`insert into orders
	// 	(user_id, amount, description)
	// 	values ($1, $2, $3)`, userID, amount, desc,
	// 	)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }

	type Order struct {
		ID, userID, amount int
		description        string
	}
	var orders []Order

	userID := 4
	rows, err := db.Query(`
	select id, amount, description
	from orders
	where user_id=$1;
	`, userID)

	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var order Order
		order.userID = userID
		err := rows.Scan(&order.ID, &order.amount, &order.description)
		if err != nil {
			panic(err)
		}
		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}

	fmt.Printf("The orders of user %d are: %v", userID, orders)
}

func main() {
	sqlExs()
}
