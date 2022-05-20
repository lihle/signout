package handler

import (
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
