package main

import (
	"database/sql"
	"webGo/models"
)

func main() {
	db := Must(models.ConnectToDB()).(*sql.DB)
	Must(nil, db.Ping())

	defer db.Close()
	userService := models.UserService{
		DB: db,
	}
	Must(userService.Create("kcat@example.com", "secret pass"))
}

func Must(value any, err error) any {
	if err != nil {
		panic(err)
	}

	return value
}
