package main

import "fmt"

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/projects")
}
