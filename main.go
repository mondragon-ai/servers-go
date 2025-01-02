package main

import (
	"fmt"
	"net/http"
)

func main () {
	fmt.Println("starting go server....")

	fileServer := http.FileServer(http.Dir("."))

	mux := http.NewServeMux()

	mux.Handle("/", fileServer)
	
	server := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}