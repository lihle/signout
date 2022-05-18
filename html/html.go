package html

import "fmt"

//H1 gives the first heading. Expects body, class, id
func H1(args ...interface{}) string {
	return general("h1", args...)
}

//H2 gives the second heading. Expects body, class, id
func H2(args ...interface{}) string {
	return general("h2", args...)
}

//H3 gives the third heading. Expects body, class, id
func H3(args ...interface{}) string {
	return general("h3", args...)
}

//H4 gives the fourth heading. Expects body, class, id
func H4(args ...interface{}) string {
	return general("h4", args...)
}

//A is a link. Expects link, body, class, id in that order
func A(args ...interface{}) string {
	var link, body, class, id string
	if len(args) > 0 {
		link = fmt.Sprint(args[0])
	}
	if len(args) > 1 {
		body = fmt.Sprint(args[1])
	}
	if len(args) > 2 {
		class = fmt.Sprint(args[2])
	}
	if len(args) > 3 {
		id = fmt.Sprint(args[3])
	}
	return "<a href='" + link + "' class='" + class + "' id='" + id + "'>" + body + "</a>"
}

//Button returns a link with a button class
func Button(href, body string, class ...string) string {
	c := "button"
	if len(class) > 0 {
		c += " " + class[0]
	}
	return Div(A(href, body, c), "button")
}

//ScriptButton triggers a javascript function
func ScriptButton(title, function string, args ...interface{}) string {
	var a string
	for i, arg := range args {
		if arg == "this" {
			a += fmt.Sprint(arg)
		} else {
			a += "'" + fmt.Sprint(arg) + "'"
		}
		if i < len(args)-1 {
			a += ","
		}
	}
	f := function + "(" + a + ")"
	return Div("<a onclick=\""+f+"\" class='button'>"+title+"</a>", "button")
}

//B is a bold tag. Expects body, class, id
func B(args ...interface{}) string {
	return general("b", args...)
}

//Br is <br/>
func Br() string {
	return "<br/>"
}

//Div is a html div. Expects body, class, id
func Div(args ...interface{}) string {
	return general("div", args...)
}

//Error is a div with the .error class
func Error(body interface{}) string {
	return Div(body, "error")
}

//Success is a div with the .success class
func Success(body string) string {
	return Div(body, "success")
}

// Img returns an image
func Img(label, src string) string {
	return `<img alt="` + label + `" src="` + src + `"/>`
}

//Label is an html label. Expects body, class, id
func Label(args ...interface{}) string {
	return general("label", args...)
}

//Span is a html span. Expects body, class, id
func Span(args ...interface{}) string {
	return general("span", args...)
}

//Tableparts

//Table holds rows. Expects body, class, id
func Table(args ...interface{}) string {
	return general("table", args...)
}

//Thead is a <thead>. Expects body, class, id
func Thead(args ...interface{}) string {
	return general("thead", args...)
}

//Tbody is a <tbody>. Expects body, class, id
func Tbody(args ...interface{}) string {
	return general("tbody", args...)
}

//Th is a table header
func Th(cells string, cols int) string {
	return fmt.Sprintf("<th colspan='%d'>%s</th>", cols, cells)
}

//Tr is a row. Expects body, class, id
func Tr(args ...interface{}) string {
	return general("tr", args...)
}

//Td is a cell. Expects body, class, id
func Td(args ...interface{}) string {
	return general("td", args...)
}

// expects body, class, id
func general(tag string, args ...interface{}) string {
	var class, val, id string
	if len(args) > 0 {
		val = fmt.Sprint(args[0])
	}
	if len(args) > 1 {
		class = fmt.Sprintf("class='%s'", args[1])
	}
	if len(args) > 2 {
		id = fmt.Sprintf(`class="%s"`, args[2])
	}
	return "<" + tag + " " + class + " " + id + ">" + val + "</" + tag + ">" //<b>body</b>
}
