package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("POST /", inputHandler)
}

func inputHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println(req)
}
