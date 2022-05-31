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
	return sql.Open("mysql", "root:Start@123@(10.1.1.2:3306)/signout?charset=utf8")
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
	_, err := db.Exec("UPDATE device_type SET quantity = ? WHERE device_type_id = ?", quantity, deviceid)
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

//AccessApproval : func -> returns 1 for correct user details & 0 for wrong details
func AccessApproval(username, password string) (loginid string, id int, err error) {
	raw := db.QueryRow("SELECT person_id , COUNT(*) FROM admin_login WHERE username LIKE ? AND password LIKE ?", username, password)
	err = raw.Scan(&loginid, &id)
	return
}

//GetAdminUser : func -> returns Admin data from username & password
func GetAdminUser(personid string) (usr types.Admin, err error) {
	raw := db.QueryRow("SELECT al.person_id, p.full_name, p.program_id, pr.program_name,al.username, al.password "+
		"FROM admin_login al LEFT JOIN person p ON al.person_id = p.person_id LEFT JOIN program pr on "+
		"pr.program_id = p.program_id WHERE al.person_id = ?", personid)
	err = raw.Scan(&usr.Person.PersonID, &usr.Person.Fullname, &usr.Program.ProgramID, &usr.Program.ProgramName,
		&usr.Username, &usr.Password)
	return
}

//InsertLogin : func -> enters a new login to histroy
func InsertLogin(personid, programid string) (string, error) {
	raw, err := db.Exec("INSERT INTO login_history SET person_id = ?, program_id = ?, logged_in = NOW()", personid,
		programid)
	if err != nil {
		return "", err
	}
	id, err := raw.LastInsertId()
	return fmt.Sprint(id), err
}

//UpdateLogoutCols : func -> updates login history when one logs in without logging out first
func UpdateLogoutCols(personid string) error {
	_, err := db.Exec("UPDATE login_history SET logged_out = NOW() WHERE person_id = ? AND logged_out IS NULL ", personid)
	return err
}

//UpdateLogout : func -> updates history when logging out
func UpdateLogout(historyid string) error {
	_, err := db.Exec("UPDATE login_history SET logged_out = NOW() WHERE login_history_id = ?", historyid)
	return err
}

//GetAllLoans : func -> returns all devices loaned out as loans
func GetAllLoans() (loans []types.Loan, err error) {
	raws, err := db.Query("SELECT dl.device_loan_id, dt.device_type_name, dl.device_loan_label, " +
		"dl.device_loan_timestamp, p.full_name FROM device_loan dl LEFT JOIN device_type dt ON " +
		"dl.device_type_id = dt.device_type_id LEFT JOIN person p ON dl.person_id = p.person_id " +
		"WHERE dl.device_loan_returntime IS NULL ORDER BY device_loan_timestamp DESC")
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

//GetDeviceLoan : func -> returns device loan details
func GetDeviceLoan(loanid string) (loan types.Loan, err error) {
	raw := db.QueryRow("SELECT dl.device_loan_id, dt.device_type_name, dl.device_loan_label, "+
		"dl.device_loan_timestamp, p.full_name FROM device_loan dl LEFT JOIN device_type dt ON "+
		"dl.device_type_id = dt.device_type_id LEFT JOIN person p ON dl.person_id = p.person_id "+
		"WHERE dl.device_loan_id = ?", loanid)
	err = raw.Scan(&loan.LoanID, &loan.Device.DeviceName, &loan.LoanLabel, &loan.LoanTimeStamp, &loan.Person.Fullname)
	return
}

//GetLoanDeviceType : func -> returns device type from a loan ID
func GetLoanDeviceType(loanid string) (device types.Device, err error) {
	raw := db.QueryRow("SELECT dt.* FROM device_loan dl LEFT JOIN device_type dt ON "+
		"dl.device_type_id = dt.device_type_id WHERE dl.device_loan_id = ?", loanid)
	err = raw.Scan(&device.DeviceID, &device.DeviceName, &device.ProgramID, &device.Quantity)
	return
}

//UpdateDeviceQuantity : func -> updates quantity when signing device in
func UpdateDeviceQuantity(quantity, deviceid string) error {
	_, err := db.Exec("UPDATE device_type SET quantity = ? WHERE device_type_id = ?", quantity, deviceid)
	return err
}

//UpdateSignin : func -> updates loan when signning in loaned out device
func UpdateSignin(userid, comment, loanid string) error {
	_, err := db.Exec("UPDATE device_loan SET device_loan_returntime = NOW(), return_user_id = ?, "+
		"device_loan_comment = ? WHERE device_loan_id = ?", userid, comment, loanid)
	return err
}

//InsertDeviceType : func -> inserts a new entry for program
func InsertDeviceType(typename, programid, quantity string) error {
	_, err := db.Exec("INSERT INTO device_type SET device_type_name = ?, program_id = ?, "+
		"quantity = ?", typename, programid, quantity)
	return err
}

//InsertProgram : func -> insert new entry for new program
func InsertProgram(name, desc string) error {
	_, err := db.Exec("INSERT INTO program SET program_name = ?, program_definition = ?", name, desc)
	return err
}

//CountSignedout : func -> returns a count of all signed out per device type id
func CountSignedout(dtypeid string) (count int, err error) {
	raw := db.QueryRow("SELECT COUNT(*) FROM device_loan dl WHERE dl.device_type_id = ? "+
		"AND dl.device_loan_returntime IS NULL", dtypeid)
	err = raw.Scan(&count)
	return
}

//GetProgram : func -> returns program details per program ID
func GetProgram(programid string) (program types.Program, err error) {
	raw := db.QueryRow("SELECT * FROM program WHERE program_id = ?", programid)
	err = raw.Scan(&program.ProgramID, &program.ProgramName, &program.ProgramDefinition)
	return
}

//UpdateProgram : func -> updates the program date
func UpdateProgram(name, define, programid string) error {
	_, err := db.Exec("UPDATE program SET program_name = ?, program_definition = ? "+
		"WHERE program_id = ?", name, define, programid)
	return err
}

//GetDevice : func -> returns device type information per id
func GetDevice(deviceid string) (device types.Device, err error) {
	raw := db.QueryRow("SELECT * FROM device_type WHERE device_type_id = ?", deviceid)
	err = raw.Scan(&device.DeviceID, &device.DeviceName, &device.ProgramID, &device.Quantity)
	return
}
