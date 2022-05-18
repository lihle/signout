package html

//LabelString creates a <input type="string"/> with wrappers. Expects label, name, value, attributes
func LabelString(label string, args ...interface{}) string {
	return Div(Label(label) + String(args...))
}

//LabelTextArea creates a <textarea name="" attr=""> val </textarea> with wrappers. Expects name, value, rows, cols
func LabelTextArea(label string, args ...interface{}) string {
	return Div(Label(label) + TextArea(args...))
}

//LabelFilePicker creates a filepicker with wrappers. Expects label, name, extensions
func LabelFilePicker(label, name, extensions string) string {
	return Div(Label(label) + FilePicker(name, extensions))
}

//LabelDate returns an <input type="date"/> with wrappers
func LabelDate(label string, args ...interface{}) string {
	return Div(Label(label) + Date(args...))
}

//LabelSelect is a <select><option></select> with labels and wrappers
//and an optional preselected item
func LabelSelect(label, name string, ids, labels []string, args ...interface{}) string {
	return Div(Label(label) + Select(name, ids, labels, args...))
}

//LabelCheckbox is a checkbox + a label
func LabelCheckbox(label, name string, checked bool) string {
	return Div(Label(label) + Checkbox(name, checked))
}
