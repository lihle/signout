package main

import (
	"fmt"
	"net/http"
	"os"
	"signout/handler"
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

	//Pages
	r.HandleFunc("/", handler.ViewHomepage)
	r.HandleFunc("/signout_device", handler.ViewSignoutpage)
	r.HandleFunc("/add_person", handler.ViewAddPerson)

	//Static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	port := ":5001"
	fmt.Println("http://localhost" + port)
	fmt.Println()
	http.ListenAndServe(port, r)
}
