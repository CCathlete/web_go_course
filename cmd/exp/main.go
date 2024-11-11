package main

import (
	"html/template"
	"os"
)

func main() {
	t, err := template.ParseFiles("hello.gohtml")
	if err != nil {
		panic(err)
	}

	user := struct{ Name, Bio string }{
		Name: "Ken",
		Bio:  "<script>alert(Haha you've been h4x0r3d!);</script>",
	}

	err = t.Execute(os.Stdout, user)
	if err != nil {
		panic(err)
	}
}
