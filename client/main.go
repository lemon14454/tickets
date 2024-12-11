package main

import (
	"fmt"
	"net/http"
)

func main() {

	static := http.Dir("web/dist")

	http.Handle("/", http.FileServer(static))

	fmt.Println("Client Started on PORT 3000 !")
	http.ListenAndServe("localhost:3000", nil)
}
