package types

//Person: class
type Person struct {
	PersonID  string
	Fullname  string
	ProgramID string
}

//Admin : class
type Admin struct {
	Person   Person
	Program  Program
	Username string
	Password string
}

//AdminHistory : class
type AdminHistory struct {
	HistoryID string
	Admin     Admin
	LoggedIn  string
	LoggedOut string
}
