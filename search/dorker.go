package search

import "strings"

/*
Dorker is an object that constructs Google dorks for searching.
*/
type Dorker struct {
	elements []string
	query    string
}

/*
NewDorker creates a new Dorker object.
*/
func NewDorker(elements []string) *Dorker {
	return &Dorker{elements: elements}
}

/*
All builds the query string and returns it without any modifications.
This is mean for a sweeping search across all file types indexed by Google.
*/
func (d *Dorker) All() string {
	d.buildQuery()
	return d.query
}

/*
Docs builds the query string and returns it with the "filetype:doc OR filetype:pdf" operator.
This is meant for searching for documents indexed by Google.
*/
func (d *Dorker) Docs() string {
	d.buildQuery()
	return d.query + " filetype:doc OR filetype:pdf"
}

/*
buildQuery builds a query string from the elements.
The most common case will be to wrap each element into an "intext:" operator, inside quotes.
*/
func (d *Dorker) buildQuery() {
	// Wrap each element in quotes and intext: operator
	for i, element := range d.elements {
		d.elements[i] = "intext:\"" + element + "\""
	}

	// Join the elements into a single query string
	d.query = strings.Join(d.elements, " ")
}