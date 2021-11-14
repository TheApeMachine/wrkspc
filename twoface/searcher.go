package twoface

/*
Searcher is a more convenient way of going through a common object we use,
the slice of key/values
*/
func Searcher(space []map[string]string, key string) string {
	for _, block := range space {
		if block[key] != "" {
			return block[key]
		}
	}

	return ""
}
