package main

import (
	"fmt"
	"net/http"
	"os"
	"signout/storage"

	"github.com/gorilla/mux"
)

//The main function
func main() {
	fmt.Println("Hello-World")
	//
	err := storage.Init()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	r := mux.NewRouter()

	//Static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	port := ":5001"
	fmt.Println("http://localhost" + port)
	fmt.Println()
	http.ListenAndServe(port, r)
}
