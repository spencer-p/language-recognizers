package dfa_test

import (
	"fmt"
	"testing"
	"testing/quick"

	"github.com/spencer-p/language-recognizers/dfa"
)

// intToBinaryString generates strings of 0*1* from integers. Useful for testing
// DFAs where the language is {0, 1}.
func intToBinaryString(i int) string {
	return fmt.Sprintf("%b", i)
}

// TestEvenOnes makes sure a DFA for an even amount of 1s in the strings works.
func TestEvenOnes(t *testing.T) {
	evenOnes := func(s string) bool {
		count := 0
		for i := 0; i < len(s); i++ {
			if s[i] == '1' {
				count += 1
			}
		}
		return count%2 == 0
	}

	// A DFA that accepts strings with even amounts of ones has two states.
	// They flop between each other on ones, and transition to themselves on any
	// other character. The start state is the accept state.

	transition := func(q int, a byte) int {
		if q == 1 && a == '1' {
			return 2
		} else if q == 2 && a == '1' {
			return 1
		} else {
			return q
		}
	}

	accept := map[int]bool{1: true}

	f := func(i int) bool {
		s := intToBinaryString(i)
		return evenOnes(s) == recognize(1, transition, accept, s)
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}
