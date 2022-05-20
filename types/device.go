package types

//Device : class
type Device struct {
	DeviceID   string
	DeviceName string
	ProgramID  string
	Quantity   int
}

//Loan : class
type Loan struct {
	LoanID        string
	Device        Device
	Person        Person
	LoanDate      string
	LoanTimeStamp string
	LoanLabel     string
}
