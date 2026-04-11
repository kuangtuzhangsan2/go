package main

import (
	"fmt"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>404 not found(doge)</h1>")
}

func main() {
	http.HandleFunc("/hello", hello)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		fmt.Println("http serve failed, err:", err)
		return
	}
	fmt.Println("Server is running on port 9090...")
}