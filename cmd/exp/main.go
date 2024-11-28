package main

import "fmt"

func main() {
	password := "this is a totally secret password nobody will guess"
	hash, err := HashHMAC(password)
	if err != nil {
		panic(err)
	}
	fmt.Println(hash)
}
