package table

import (
	"fmt"
	"signout/html"
)

// Table holds the table data
type Table struct {
	headings []header
	rows     [][]string
}

type header struct {
	name string
	cols int
}

// New returns a new table
func New(headings ...interface{}) *Table {
	var table Table
	for _, head := range headings {
		table.headings = append(table.headings, header{fmt.Sprint(head), 1})
	}
	return &table
}

// AddHeading adds a heading to the table
func (t *Table) AddHeading(head interface{}, cols ...int) {
	if len(cols) == 0 {
		t.headings = append(t.headings, header{fmt.Sprint(head), 1})
	} else {
		t.headings = append(t.headings, header{fmt.Sprint(head), cols[0]})
	}
}

// AddRow adds a row to the table
func (t *Table) AddRow(cells ...interface{}) {
	var scells []string
	for _, cell := range cells {
		scells = append(scells, fmt.Sprint(cell))
	}
	t.rows = append(t.rows, scells)
}

// AddCell adds a cell to the most recent row
func (t *Table) AddCell(cell interface{}) {
	t.rows[len(t.rows)-1] = append(t.rows[len(t.rows)-1], fmt.Sprint(cell))
}

// HTML returns the html of a table
func (t *Table) HTML(class string) string {
	var thead string
	for _, heading := range t.headings {
		thead += html.Th(heading.name, heading.cols)
	}
	thead = html.Thead(html.Tr(thead))

	var tbody string
	for _, row := range t.rows {
		var r string
		for _, cell := range row {
			r += html.Td(cell)
		}
		tbody += html.Tr(r)
	}
	tbody = html.Tbody(tbody)
	return html.Table(thead+tbody, class)
}
