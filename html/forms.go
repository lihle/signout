package html

import "fmt"

//This package generates form data

//Form returns the body wrapped in a <form> tag
func Form(body, method, action string) string {
	return "<form method='" + method + "' action='" + action + "'>" + body + "</form>"
}

//Submit is a submit button
func Submit(val string) string {
	return input("submit", "submit", val, "")
}

//Date is a form date. Expects name, value, attributes
func Date(args ...interface{}) string {
	return input("date", args...)
}

//FilePicker is a filepicker
func FilePicker(name, extensions string) string {
	extensions = "accept='" + extensions + "'"
	return input("file", name, "", extensions)
}

//Hidden produces an hidden input. Expects name, value, attributes
func Hidden(args ...interface{}) string {
	return input("hidden", args...)
}

//Select Generates a <select><option> with a preselected choice
func Select(name string, ids, labels []string, args ...interface{}) string {
	var id string
	if len(args) > 0 {
		id = fmt.Sprint(args[0])
	}
	s := "<select name='" + name + "' id='" + name + "'>"
	for i := range labels {
		s += option(ids[i], labels[i], ids[i] == id)
	}
	s += "</select>"
	return s
}

//String is a form input type=string. Expects name, value, attributes
func String(args ...interface{}) string {
	return input("string", args...)
}

//TextArea is a form text. Expects name, value, rows, cols
func TextArea(args ...interface{}) string {
	var name, value string
	rows := "4"
	cols := "50"
	if len(args) > 0 {
		name = fmt.Sprint(args[0])
	}
	if len(args) > 1 {
		value = fmt.Sprint(args[1])
	}
	if len(args) > 2 {
		rows = fmt.Sprint(args[2])
	}
	if len(args) > 3 {
		cols = fmt.Sprint(args[3])
	}
	return `<textarea name="` + name + `" rows="` + rows + `" cols="` + cols + `">` + value + `</textarea>`
}

//ScriptInput triggers a javascript function on change
func ScriptInput(title, function string, args ...string) string {
	var a string
	for i, arg := range args {
		if arg == "this" {
			a += arg
		} else {
			a += "'" + arg + "'"
		}
		if i < len(args)-1 {
			a += ","
		}
	}
	f := function + "(" + a + ")"
	return fmt.Sprintf(`<label>%s</label><input type="string" onInput="%s"/>`, title, f)
}

//Checkbox is an active checkbox.
func Checkbox(name string, checked bool) string {
	var c string
	if checked {
		c = "checked"
	}
	return input("checkbox", name, "", c)
}

//Custombox is a custom checkbox built with images
func Custombox(name string, checked bool, color, fn string) string {
	var c string
	if checked {
		c = "checked"
	}
	return "<div class='custombox " + color + " " + c + //Make it look right
		"' onClick='$(this).toggleClass(\"checked\");" + //Make it toggle the check
		fn + "(this,\"" + name + "\")'></div>" // Do something custom
}

//DisabledCheckbox is a prop
func DisabledCheckbox(checked bool) string {
	attr := "disabled "
	if checked {
		attr += "checked"
	}
	return input("checkbox", "", "", attr)
}

//Expects tag, name, value, attributes
func input(tag string, args ...interface{}) string {
	var name, val, attr string
	if len(args) > 0 {
		name = fmt.Sprintf(`name="%s" id="%s"`, args[0], args[0])
	}
	if len(args) > 1 {
		val = fmt.Sprintf(`value="%s"`, args[1])
	}
	if len(args) > 2 {
		attr = fmt.Sprint(args[2])
	}
	return "<input type='" + tag + "' " + name + " " + val + " " + attr + "/>"
}

func option(name, label string, selected bool) string {
	var s string
	if selected {
		s = "selected"
	}
	return "<option value='" + name + "' " + s + " >" + label + "</option>"
}

func datalist(id string, ids, labels []string) string {
	total := "<datalist id='" + id + "'>"
	for i := range ids {
		total += option(ids[i], labels[i], false)
	}
	total += "</datalist>"
	return total
}

//Autocomplete returns an input and datalist matching awesompletes prefs
func Autocomplete(name string, ids, labels []string) string {
	id := name + "_list"
	total := "<input name='student' class='awesomplete dropdown-input' list='" + id + "' />"
	total += datalist(id, ids, labels)
	return total
}
