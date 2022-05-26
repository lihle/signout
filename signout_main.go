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
	r.HandleFunc("/admin_login", handler.ViewAdminlogin)
	r.HandleFunc("/admin_verify", handler.AdminLoginVerify)
	r.HandleFunc("/admin_user", handler.ViewAdminUser)
	r.HandleFunc("/admin_logout", handler.AdminLogout)

	//
	r.HandleFunc("/admin_user/device_loan", handler.ViewAdminLoanSignin)
	r.HandleFunc("/admin_user/persons", handler.ViewAllPersons)
	r.HandleFunc("/admin_user/devices", handler.ViewAllDevices)
	r.HandleFunc("/admin_user/add_device_type", handler.ViewAddNewDevice)
	r.HandleFunc("/admin_user/add_new_program", handler.ViewAddNewProgram)
	r.HandleFunc("/admin_user/extra/program", handler.EditProgramDetails)
	r.HandleFunc("/admin_user/loanout", handler.ViewAdminLoanout)

	//Static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	port := ":5001"
	fmt.Println("http://localhost" + port)
	fmt.Println()
	http.ListenAndServe(port, r)
}
