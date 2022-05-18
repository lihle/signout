package handler

import (
	"fmt"
	"html/template"
	"net/http"
)

// Contains Page struct, functions, template parsing and utility functions

const tmpl = "template.html"

var templates = template.Must(template.ParseFiles(tmpl))

//Page is what is exported to the template
type Page struct {
	Title string
	Body  template.HTML //string
	Error template.HTML //string
	Form  Form
}

//Form keeps the form data
type Form struct {
	Method string
	Action string
	Enc    string
	Body   template.HTML
}

func view(w http.ResponseWriter, page Page) {
	templates.ExecuteTemplate(w, tmpl, page)
}

func newPage(title, body string) Page {
	return Page{Title: title, Body: template.HTML(body)}
}

func newForm(title, form string) Page {
	return newFormExtended(title, form, "", "", "application/x-www-form-urlencoded")
}

func multiPartForm(title, form string) Page {
	return newFormExtended(title, form, "", "post", "multipart/form-data")
}

func newFormExtended(title, form, action, method, enc string) Page {
	f := Form{Body: template.HTML(form), Method: "post"}
	f.Enc = enc
	f.Action = action
	if len(method) > 0 {
		f.Method = method
	}
	return Page{Title: title, Form: f}
}

//SetBody sets the value of the body
func (page *Page) SetBody(content string) {
	page.Body = template.HTML(content)
}

//SetForm sets the value of the form
func (page *Page) SetForm(content string) {
	tmp := newForm("", content)
	page.Form = tmp.Form
}

//SetError sets the errorMessage
func (page *Page) SetError(content string) {
	page.Error = template.HTML(content)
}

func isErr(w http.ResponseWriter, r *http.Request, err error) bool {
	if err != nil {
		fmt.Println(err)
		templates.ExecuteTemplate(w, tmpl, newPage("Error", err.Error()))
		// http.Redirect(w, r, "/404", http.StatusSeeOther)
		return true
	}
	return false
}

func set(vals ...string) bool {
	for _, s := range vals {
		if len(s) == 0 {
			return false
		}
	}
	return true
}
