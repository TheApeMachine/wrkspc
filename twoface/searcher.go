package twoface

import "github.com/theapemachine/wrkspc/spdg"

/*
Searcher is a more convenient way of going through a common object we use,
the slice of key/values
*/
func Searcher(space []spdg.Annotation, key string) string {
	for _, block := range space {
		if block.Key != "" {
			return block.Key
		}
	}

	return ""
}
