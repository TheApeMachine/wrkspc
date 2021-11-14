package luthor

func isOpeningDelimiter(token string) (bool, int) {
	for _, pair := range delimiters {
		if ok1, _, n := checkPair(pair, string(getLastChar(token, 1))); ok1 {
			return ok1, n
		}
	}

	return false, 0
}

func inDelimiter(token string) (bool, int) {
	for _, pair := range delimiters {
		if _, ok2, n := checkPair(pair, token); ok2 {
			return ok2, n
		}
	}

	return false, 0
}

func tokenBetweenDelimiters(token string, n int) string {
	max := Max(1, len(token)-n)

	inner := token[1:max]
	return inner
}

func getLastChar(token string, n int) []rune {
	return []rune(token)[len(token)-n : len(token)]
}

func getFirstChar(token string, n int) []rune {
	return []rune(token)[0:n]
}

func getLastString(token string, n int) string {
	return string(getLastChar(token, n))
}

func checkPair(pair []string, token string) (bool, bool, int) {
	var ok1, ok2 bool
	var pos int

	if pair[0] == string(getFirstChar(token, 1)) {
		ok1 = true
		pos = 1
	}

	if len(token) >= 2 && pair[0] == string(getFirstChar(token, 2)) {
		ok1 = true
		pos = 2
	}

	if pair[1] == string(getLastChar(token, 1)) {
		ok2 = true
		pos = 1
	}

	if len(token) >= 4 && pair[1] == string(getLastChar(token, 2)) {
		ok2 = true
		pos = 2
	}

	return ok1, ok2, pos
}
