package luthor

/*
Entity is a struct that represents a single entity.
*/
type Entity struct {
	Name  string  `json:"entity"`
	Score float64 `json:"score"`
	Index int     `json:"index"`
	Word  string  `json:"word"`
	Start int     `json:"start"`
	End   int     `json:"end"`
}
