package forms

import (
	"net/http"
	"signout/html"
	"signout/storage"
)

//StudentSelect returns a select with all students
func StudentSelect(r *http.Request) (string, error) {
	name := "student"
	ids, labels, err := studentsIntoIdsAndLabels()
	if err != nil {
		return "", err
	}

	return html.Select(name, ids, labels, r.FormValue(name)), nil
}

//StudentAutocomplete is a search function using awesomplete
func StudentAutocomplete() (string, error) {
	ids, labels, err := studentsIntoIdsAndLabels()
	if err != nil {
		return "", err
	}
	total := html.Autocomplete("student", ids, labels)
	return html.Div(html.Label("Person (autocomplete):") + total), nil
}

func studentsIntoIdsAndLabels() (ids, labels []string, err error) {
	persons, err := storage.GetAllPersons()
	if err != nil {
		return
	}
	for _, person := range persons {
		ids = append(ids, person.PersonID)
		labels = append(labels, person.Fullname)
	}
	return
}
