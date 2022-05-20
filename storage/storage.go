package storage

import (
	"database/sql"
	"fmt"
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

//GetAllPrograms() : func
func GetAllPrograms() (programs []types.Program, err error) {
	raws, err := db.Query("SELECT * FROM program")
	if err != nil {
		return
	}
	defer raws.Close()

	for raws.Next() {
		var program types.Program
		err = raws.Scan(&program.ProgramID, &program.ProgramName, &program.ProgramDefinition)
		if err != nil {
			return
		}
		programs = append(programs, program)
	}
	return
}

//GetAllDevices : func -> returns all device types, their quantity as per program
func GetAllDevices(programid string) (devices []types.Device, err error) {
	raws, err := db.Query("SELECT * FROM device_type WHERE quantity > 0 AND program_id = ?", programid)
	if err != nil {
		return
	}
	defer raws.Close()

	for raws.Next() {
		var device types.Device
		err = raws.Scan(&device.DeviceID, &device.DeviceName, &device.ProgramID, &device.Quantity)
		if err != nil {
			return
		}
		devices = append(devices, device)
	}
	return
}

//GetPerson : func -> returns person using person id
func GetPerson(id string) (person types.Person, err error) {
	raw := db.QueryRow("SELECT * FROM person WHERE person_id = ?", id)
	err = raw.Scan(&person.PersonID, &person.Fullname, &person.ProgramID)
	return
}

//GetPogramID : func -> returns program id from program name
func GetPogramID(name string) (id string, err error) {
	raw := db.QueryRow("SELECT program_id FROM program WHERE program_name = ?", name)
	err = raw.Scan(&id)
	return
}

//InsertPerson : func -> add person data and returns person id
func InsertPerson(name, programID string) (string, error) {
	raw, err := db.Exec("INSERT INTO person set full_name = ?, program_id = ?", name, programID)
	if err != nil {
		return "", err
	}
	id, err := raw.LastInsertId()
	return fmt.Sprint(id), err
}

//GetDeviceType : func -> returns device type info from type name & program id
func GetDeviceType(name, programID string) (device types.Device, err error) {
	raw := db.QueryRow("SELECT * FROM device_type WHERE device_type_name = ? and program_id = ?", name, programID)
	err = raw.Scan(&device.DeviceID, &device.DeviceName, &device.ProgramID, &device.Quantity)
	return
}

//InsertDeviceLoan : func -> add device loan and returns an error if it exists
func InsertDeviceLoan(typeid, personid, label string) error {
	_, err := db.Exec("INSERT INTO device_loan SET device_type_id = ?, person_id = ?, device_loan_date = NOW(), "+
		"device_loan_timestamp = NOW(), device_loan_label = ?", typeid, personid, label)
	if err != nil {
		return err
	}
	return err
}

//UpdateTypeQuantity : func -> updates device type quantity and returns an error if it exists
func UpdateTypeQuantity(quantity int, deviceid string) error {
	_, err := db.Exec("UPDATE device_type set quantity = ? WHERE device_type_id = ?", quantity, deviceid)
	return err
}

//GetLoanedOut : func -> returns all loaned out program devices per program id
func GetLoanedOut(programID string) (loans []types.Loan, err error) {
	raws, err := db.Query("SELECT dl.device_loan_id, dt.device_type_name, dl.device_loan_label, "+
		"dl.device_loan_timestamp, p.full_name FROM device_loan dl LEFT JOIN device_type dt ON "+
		"dl.device_type_id = dt.device_type_id LEFT JOIN person p ON dl.person_id = p.person_id "+
		"WHERE dt.program_id = ? and dl.device_loan_returntime IS NULL ", programID)
	if err != nil {
		return
	}
	defer raws.Close()

	for raws.Next() {
		var loan types.Loan
		err = raws.Scan(&loan.LoanID, &loan.Device.DeviceName, &loan.LoanLabel, &loan.LoanTimeStamp, &loan.Person.Fullname)
		if err != nil {
			return
		}
		loans = append(loans, loan)
	}
	return
}
