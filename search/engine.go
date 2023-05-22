package search

import (
	customsearch "google.golang.org/api/customsearch/v1"
	"google.golang.org/api/googleapi/transport"
)

/*
Engine is a wrapper around the Google Custom Search Engine API.
*/
type Engine struct {
}

/*
NewEngine creates a new Engine object.
*/
func NewEngine() *Engine {
	return &Engine{}
}

/*
Query performs a search against the Google Custom Search Engine API.
*/
func (e *Engine) Search(q string) {
	
}