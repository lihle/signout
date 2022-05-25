package handler

import (
	"net/http"
	"signout/storage"
)

//ProgramList : func -> returns an array type string of a list of available programs
func ProgramList() (programs []string) {
	programs = append(programs, "")
	pgs, err := storage.GetAllPrograms()
	if err != nil {
		return
	}

	for _, p := range pgs {
		programs = append(programs, p.ProgramName)
	}
	return
}

//DeviceList : func -> returns an array type string of a list of available devices
func DeviceList(programid string) (devices []string) {
	devices = append(devices, "")
	dvs, err := storage.GetAllDevices(programid)
	if err != nil {
		return
	}
	for _, d := range dvs {
		devices = append(devices, d.DeviceName)
	}
	return
}

//AdminLoginVerify : func -> it simply verify the provided login details and gives admin access
func AdminLoginVerify(w http.ResponseWriter, r *http.Request) {
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	//check returned value
	access, err := storage.AccessApproval(username, password)
	if err == nil {
		if access == 1 {
			user, err := storage.GetAdminUser(username, password)
			if err == nil {
				err := storage.UpdateLogoutCols(user.PersonID)
				if err != nil {
					return
				}
				id, err := storage.InsertLogin(user.PersonID, user.Fullname, user.ProgramID, user.ProgramName)
				if err == nil {
					http.Redirect(w, r, "/admin_user?u="+id, http.StatusSeeOther)
				}
			}
		} else if access == 0 {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}
}
