package handler

import (
	"net/http"
	"signout/html"
	"signout/html/forms"
	"signout/html/table"
	"signout/storage"
)

//ViewHomepage : func -> home page of the system
func ViewHomepage(w http.ResponseWriter, r *http.Request) {
	f, err := forms.StudentAutocomplete()
	if err != nil {
		return
	}
	//
	body := html.Div(html.A("/admin_login", "[login]"), "right")
	body += html.Br()
	body += html.H2("type your name to search if exist, select it then press 'Submit'")
	body += f
	body += html.Br()
	body += html.H2("else if doesn't exist, follow link below to enter your name & surname")
	body += html.Br()
	body += html.A("/add_person", "(to add your name)")
	//
	if set(r.FormValue("submit")) {
		u := r.FormValue("student")
		http.Redirect(w, r, "/signout_device?u="+u, http.StatusSeeOther)
	}

	view(w, multiPartForm("Tech signout system", body))
}

//ViewSignoutpage : func -> the page for actually signing out a device
func ViewSignoutpage(w http.ResponseWriter, r *http.Request) {
	var body string
	u := r.FormValue("u")

	p, _ := storage.GetPerson(u) //get person info
	devices := DeviceList(p.ProgramID)
	loans, err := storage.GetLoanedOut(p.ProgramID)
	if err != nil {
		return
	}

	//
	body += html.A("/", "(Go-to home)")
	body += html.H2(" ")
	body += html.LabelSelect("Select device (only the devices, one has access to _) : ", "device", devices, devices)
	body += html.Br()
	body += html.LabelString("Label/Comment (if labelled, write label) : ", "label", "No label")
	body += html.Br()
	body += html.H2("All signed out device", "clickable")
	//
	t := table.New("#No", "Device (type of the device)", "Label (label on device)", "Date (signed out at)",
		"Person (signed out by)")
	for x, loan := range loans {
		t.AddRow(x+1, loan.Device.DeviceName, loan.LoanLabel, loan.LoanTimeStamp, loan.Person.Fullname)
	}
	body += html.Div(t.HTML("tablesorter"), "hidden")

	if set(r.FormValue("submit")) {
		label := r.FormValue("label")
		devicename := r.FormValue("device")
		device, err := storage.GetDeviceType(devicename, p.ProgramID)
		if err != nil {
			return
		}
		err = storage.InsertDeviceLoan(device.DeviceID, u, label)
		if err == nil {
			device.Quantity = device.Quantity - 1
			err = storage.UpdateTypeQuantity(device.Quantity, device.DeviceID)
			if err == nil {
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}
		}
	}

	view(w, multiPartForm(p.Fullname+" : Signing out a device", body))
}

//ViewAddPerson : func -> the page for adding a person's name
func ViewAddPerson(w http.ResponseWriter, r *http.Request) {
	var body string
	programs := ProgramList()

	//
	body += html.A("/", "(Go to home)")
	body += html.H2("Please enter your full name and surname below..")
	body += html.Br()
	body += html.LabelString("First-Name : ", "name")
	body += html.Br()
	body += html.LabelString("Last-Name : ", "surname")
	body += html.Br()
	body += html.H2("of the Axium programs, which are you linked with?")
	body += html.Br()
	body += html.LabelSelect("Program (link with Axium) : ", "program", programs, programs)

	//
	if set(r.FormValue("submit")) {
		if set(r.FormValue("name"), r.FormValue("surname"), r.FormValue("program")) {
			fullname := r.FormValue("name") + " " + r.FormValue("surname")
			program := r.FormValue("program")
			pid, err := storage.GetPogramID(program)
			if err == nil {
				id, err := storage.InsertPerson(fullname, pid)
				if err != nil {
					return
				}
				http.Redirect(w, r, "/signout_device?u="+id, http.StatusSeeOther)
			}
		}
	}

	view(w, multiPartForm("Add yourself", body))
}

//ViewAdminlogin : func -> gets the login page
func ViewAdminlogin(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "login/index.html")
}

//ViewAdminUser : func -> view page for approved user
func ViewAdminUser(w http.ResponseWriter, r *http.Request) {
	var body string
	u := r.FormValue("u") //Login History id so is easy to update history table
	loans, err := storage.GetAllLoans()
	if err != nil {
		return
	}
	t := table.New("#No:", "Device (type of the device)", "Label (label on device)", "Date (signed out at)",
		"Person (signed out by)")

	body += html.Div(html.A("/admin_logout?u="+u, "(log-out to home page)"), "right")
	body += html.Br()
	body += html.H2("Basic options")
	body += html.Button("/admin_user/devices", "All options related: All Devices")
	body += html.Br()
	body += html.Button("/admin_user/persons", "All options related: Persons")
	body += html.Br()
	body += html.H2("Sign-in devices signed out")
	//
	for x, loan := range loans {
		link := "/admin_user/device_loan?u=" + u + "&loanid=" + loan.LoanID
		t.AddRow(x+1, loan.Device.DeviceName, loan.LoanLabel, loan.LoanTimeStamp, html.A(link, loan.Person.Fullname))
	}
	body += html.Div(t.HTML("tablesorter"))

	view(w, newPage("Admin User", body))
}

//AdminLogout : func -> updates the login history table and log user out
func AdminLogout(w http.ResponseWriter, r *http.Request) {
	u := r.FormValue("u")
	err := storage.UpdateLogout(u)
	if err == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

//ViewAdminLoanSignin : func -> shows device loan details and allows one to sign it in
func ViewAdminLoanSignin(w http.ResponseWriter, r *http.Request) {
	var body string
	loanid := r.FormValue("loanid")
	u := r.FormValue("u")
	loan, err := storage.GetDeviceLoan(loanid)
	if err != nil {
		return
	}
	//
	body += html.Div(html.A("/admin_user?u="+u, "(Go-to Admin home page)"), "right")
	body += html.H2("Loan details:")
	body += html.Br()
	body += html.Div(html.B("Device type (what is the device?) : ")+loan.Device.DeviceName, "overview")
	body += html.Div(html.B("Device label (written on the device) : ")+loan.LoanLabel, "overview")
	body += html.Div(html.B("Timestamp (when was it signed out?) : ")+loan.LoanTimeStamp, "overview")
	body += html.Div(html.B("Person (who loaned it out?) : ")+loan.Person.Fullname, "overview")
	body += html.LabelTextArea("Device condition : ", "comment")
	body += html.Br()
	body += html.H2("Submit to sign it in")

	if set(r.FormValue("submit")) {
		comment := "NB: " + r.FormValue("comment")
		err := storage.UpdateSignin(u, comment, loan.LoanID)
		if err != nil {
			return
		}

		device, err := storage.GetLoanDeviceType(loanid)
		if err == nil {
			device.Quantity = device.Quantity + 1
			err = storage.UpdateTypeQuantity(device.Quantity, device.DeviceID)
			if err != nil {
				return
			}
			http.Redirect(w, r, "/admin_user?u="+u, http.StatusSeeOther)
		}
	}

	view(w, multiPartForm("Sign it in", body))
}
