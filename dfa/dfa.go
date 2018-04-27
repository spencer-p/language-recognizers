package dfa

// recognize returns true if the string is accepted by the DFA represented by
// the arguments. It is essentially the extended transition function.
func recognize(start int, transition func(q int, a byte) int,
	accept map[int]bool, str string) bool {

	if len(str) == 0 {
		return accept[start]
	}

	return recognize(transition(start, str[0]), transition, accept, str[1:])
}
