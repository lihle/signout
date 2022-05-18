package storage

import (
	"database/sql"
	"signout/types"

	//this driver is needed to run sql queries to mysql

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

//Init connects to the database and returns an error upon failure
func Init() error {
	var err error
	db, err = conn()
	return err
}

func conn() (*sql.DB, error) {
	return sql.Open("mysql", "lihle:lihle@(localhost:3306)/signout?charset=utf8")
}

func valueOrWildcard(in string) string {
	if in == "" {
		return "%"
	}
	return in
}

//GetAllPersons() : func
func GetAllPersons() (persons []types.Person, err error) {
	raws, err := db.Query("SELECT * FROM person")
	if err != nil {
		return
	}
	defer raws.Close()

	for raws.Next() {
		var person types.Person
		err = raws.Scan(&person.PersonID, &person.Fullname, &person.ProgramID)
		if err != nil {
			return
		}
		persons = append(persons, person)
	}
	return
}
