package luthor

import (
	"context"

	"github.com/chromedp/chromedp"
	"golang.org/x/net/html"
)

/*
PlanType is an enum for the different types of plans that can be made.
*/
type PlanType uint

const (
	// SEARCHWIDE indicates we have a target, but no specific clues yet.
	SEARCHWIDE PlanType = iota
	// SEARCHNARROW indicates we have a target and some clues.
	SEARCHNARROW PlanType = iota
	// SEARCHDEEP indicates we have a target and a lot of clues.
	SEARCHDEEP PlanType = iota
)

/*
Planner is the top level object that makes decisions about what to do next.
*/
type Planner struct {
	buckets map[html.NodeType][]*html.Node
}

/*
NewPlanner creates a new Planner.
*/
func NewPlanner(elements []html.Node) *Planner {
	buckets := make(map[html.NodeType][]*html.Node)

	// Divide the elements into buckets, based on their HTML elemet type.
	for _, element := range elements {
		buckets[element.Type] = append(buckets[element.Type], &element)
	}

	return &Planner{buckets}
}

/*
Next returns the next action to take.
*/
func (p *Planner) Next(ctx context.Context) *chromedp.ActionFunc {
	return nil
}
