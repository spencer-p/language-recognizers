package dfa

// Recognize returns true if the string is accepted by the DFA represented by
// the arguments. It is essentially the extended transition function.
func Recognize(start int, transition func(q int, a byte) int,
	accept map[int]bool, str string) bool {

	if len(str) == 0 {
		return accept[start]
	}

	return Recognize(transition(start, str[0]), transition, accept, str[1:])
}
