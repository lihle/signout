package handler

import (
	"net/http"
	"signout/html"
	"signout/html/forms"
	"signout/html/table"
	"signout/storage"
	"strconv"
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
	body += html.Button("/admin_user/devices?u="+u, "All options related: All Devices")
	body += html.Br()
	body += html.Button("/admin_user/persons?u="+u, "All options related: Persons")
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

//ViewAllDevices : func -> views all options related to devices
func ViewAllDevices(w http.ResponseWriter, r *http.Request) {
	var body string
	u := r.FormValue("u")
	programs, err := storage.GetAllPrograms()
	if err != nil {
		return
	}

	body += html.Div(html.A("/admin_user?u="+u, "(Go-to Admin home page)"), "right")
	body += html.H2("Basic options")
	body += html.Button("/admin_user/add_device_type?u="+u, "Add new device type")
	body += html.Button("/admin_user/add_new_program?u="+u, "Add new program")
	body += html.H2("View or Edit devices per program")
	for _, program := range programs {
		body += html.H3(html.A("/admin_user/extra/program?u="+u+"&pid="+program.ProgramID, program.ProgramName), "clickable")
		devices, _ := storage.GetAllDevices(program.ProgramID)
		t := table.New("#No:", "device (What type of device?)", "Quantity (How many devices)", "Loaned out")
		for x, device := range devices {
			count, err := storage.CountSignedout(device.DeviceID)
			if err != nil {
				return
			}
			t.AddRow(x+1, html.A("/admin_user/loanout?u="+u+"&dtype="+device.DeviceID, device.DeviceName),
				device.Quantity, count)
		}
		body += html.Div(t.HTML("tablesorter"), "hidden")
	}

	view(w, newPage("All - Devices", body))
}

//ViewAdminLoanout : func -> for admin to loan out devices
func ViewAdminLoanout(w http.ResponseWriter, r *http.Request) {
	var body string
	u := r.FormValue("u")
	dtypeid := r.FormValue("dtype")
	device, _ := storage.GetDevice(dtypeid)
	programs := ProgramList()
	//
	f, err := forms.StudentAutocomplete()
	if err != nil {
		return
	}

	body += html.Div(html.A("/admin_user/devices?u="+u, "(Go back)"), "right")
	body += html.H2("Type name below to search if exist")
	body += f
	body += html.H2("If the user doesn't exist, type in first & last name & choose relavent relation")
	body += html.Div(html.LabelString("First-name : ", "name"))
	body += html.Div(html.LabelString("Last-name : ", "surname"))
	body += html.Div(html.LabelSelect("Relation to Axium : ", "relation", programs, programs))
	body += html.H2("To the above mentioned, you signing out")
	body += html.Div(html.B("Device type : "+device.DeviceName), "overview")
	body += html.Div(html.LabelString("Label on device : ", "label"), "overview")
	body += html.H2("Submit to complete sign out")

	if set(r.FormValue("submit")) {
		device.Quantity = device.Quantity - 1

		if set(r.FormValue("student")) {
			err = storage.InsertDeviceLoan(device.DeviceID, r.FormValue("student"), r.FormValue("label"))
			if err != nil {
				return
			}

		} else if set(r.FormValue("name"), r.FormValue("surname"), r.FormValue("relation")) {
			fullname := r.FormValue("name") + " " + r.FormValue("surname")
			pid, _ := storage.GetPogramID(r.FormValue("relation"))
			id, _ := storage.InsertPerson(fullname, pid)

			err = storage.InsertDeviceLoan(device.DeviceID, id, r.FormValue("label"))
			if err != nil {
				return
			}
		}
		err = storage.UpdateDeviceQuantity(strconv.Itoa(device.Quantity), device.DeviceID)
		if err != nil {
			return
		}
		http.Redirect(w, r, "/admin_user/devices?u="+u, http.StatusSeeOther)

	}

	view(w, multiPartForm("Sign out device", body))
}

//EditProgramDetails : func -> view showing all program data and can be edited
func EditProgramDetails(w http.ResponseWriter, r *http.Request) {
	var body string
	u := r.FormValue("u")
	pid := r.FormValue("pid")
	program, _ := storage.GetProgram(pid)
	devices, _ := storage.GetAllDevices(program.ProgramID)

	body += html.Div(html.A("/admin_user/devices?u="+u, "(Go back)"), "right")
	body += html.H2("Program details")
	body += html.Div(html.LabelString("Program name : ", "program", program.ProgramName))
	body += html.Div(html.LabelTextArea("Purpose (brief definition about program) : ", "purpose", program.ProgramDefinition))
	body += html.H2("Devices associate with " + program.ProgramName + " program")
	for _, device := range devices {
		count, _ := storage.CountSignedout(device.DeviceID)
		available := device.Quantity + count
		quan := strconv.Itoa(available)
		body += html.Div(html.LabelString(html.B(device.DeviceName+" (current quantity) : "), "current", quan))
		body += html.Div(html.LabelString(html.B("Add unit(s) (How many units you adding) : "), device.DeviceName, "0"))
	}

	//
	if set(r.FormValue("submit")) {
		program.ProgramName = r.FormValue("program")
		program.ProgramDefinition = r.FormValue("purpose")
		err := storage.UpdateProgram(program.ProgramName, program.ProgramDefinition, program.ProgramID)
		if err != nil {
			return
		}
		for _, d := range devices {
			//
			unit, _ := strconv.Atoi(r.FormValue(d.DeviceName))
			d.Quantity = d.Quantity + unit
			err := storage.UpdateTypeQuantity(d.Quantity, d.DeviceID)
			if err != nil {
				return
			}
		}
		http.Redirect(w, r, "/admin_user/devices?u="+u, http.StatusSeeOther)
	}

	view(w, multiPartForm("Edit Program Details", body))
}

//ViewAddNewDevice : func -> for adding a new device type
func ViewAddNewDevice(w http.ResponseWriter, r *http.Request) {
	var body string
	u := r.FormValue("u")
	programs := ProgramList()

	body += html.Div(html.A("/admin_user/devices?u="+u, "(Go to devices & programs)"), "right")
	body += html.H2("Device details")
	body += html.Div(html.LabelString("Device type (What type of device?) : ", "type"))
	body += html.Div(html.LabelSelect("Program (Which program it belongs to) : ", "program", programs, programs))
	body += html.Div(html.LabelString("Quantity (How many units?) : ", "quantity"))

	if set(r.FormValue("submit")) {
		dtype := r.FormValue("type") // devide type
		program := r.FormValue("program")
		quantity := r.FormValue("quantity")
		programid, _ := storage.GetPogramID(program)

		err := storage.InsertDeviceType(dtype, programid, quantity)
		if err != nil {
			return
		}
		http.Redirect(w, r, "/admin_user/devices?u="+u, http.StatusSeeOther)
	}

	view(w, multiPartForm("Add new device type", body))
}

//ViewAddNewProgram : func -> for adding a new program
func ViewAddNewProgram(w http.ResponseWriter, r *http.Request) {
	var body string
	u := r.FormValue("u")
	//
	body += html.Div(html.A("/admin_user/devices?u="+u, "(Go to devices & programs)"), "right")
	body += html.H2("Program details")
	body += html.Div(html.LabelString("Program (name your program) : ", "program"))
	body += html.Div(html.LabelTextArea("Description (purpose of the program) : ", "purpose"))

	if set(r.FormValue("submit")) {
		program := r.FormValue("program")
		purpose := r.FormValue("purpose")
		err := storage.InsertProgram(program, purpose)
		if err != nil {
			return
		}

		http.Redirect(w, r, "/admin_user/devices?u="+u, http.StatusSeeOther)
	}

	view(w, multiPartForm("Add new program", body))
}

//ViewAllPersons : func -> views all options related to persons
func ViewAllPersons(w http.ResponseWriter, r *http.Request) {
	var body string
	u := r.FormValue("u")

	body += html.Div(html.A("/admin_user?u="+u, "(Go-to Admin home page)"), "right")

	view(w, newPage("All - Persons", body))
}
