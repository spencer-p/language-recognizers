package dfa_test

import (
	"fmt"
	"testing"
	"testing/quick"

	"github.com/spencer-p/language-recognizers/dfa"
)

// intToBinaryString generates strings of 0*1* from integers. Useful for testing
// DFAs where the language is {0, 1}.
func intToBinaryString(i uint16) string {
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

	f := func(i uint16) bool {
		s := intToBinaryString(i)
		return evenOnes(s) == dfa.Recognize(1, transition, accept, s)
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestNo1PairsWithOddLengthBetween(t *testing.T) {
	// Inefficient method to check the condition directly
	// This is n^2! Yikes. But it's definitely correct, and we only need to it
	// to check the correctness of the linear DFA solution.
	direct := func(s string) bool {
		for i := 0; i < len(s)-1; i++ {
			if s[i] == '1' {
				for j := i + 1; j < len(s); j++ {
					if s[j] == '1' && (j-i)%2 == 0 {
						return false
					}
				}
			}
		}
		return true
	}

	// Now we define the DFA's transition func and accept states.
	transition := func(q int, a byte) int {
		switch q {
		case 1:
			switch a {
			case '0':
				return 1
			case '1':
				return 2
			}
		case 2:
			switch a {
			case '0':
				return 3
			case '1':
				return 4
			}
		case 3:
			switch a {
			case '0':
				return 2
			case '1':
				return 5
			}
		case 4:
			switch a {
			case '0':
				return 4
			case '1':
				return 5
			}
		case 5:
			return 5
		}

		// Should never happen.
		return q
	}
	accept := map[int]bool{1: true, 2: true, 3: true, 4: true}

	// Check specific examples
	if dfa.Recognize(1, transition, accept, "101") == true {
		t.Error("101 wrongly classified as true, expected false")
	}

	if dfa.Recognize(1, transition, accept, "111") == true {
		t.Error("111 wrongly classified as true, expected false")
	}

	if dfa.Recognize(1, transition, accept, "11") == false {
		t.Error("11 wrongly classified as false, expected true")
	}

	if dfa.Recognize(1, transition, accept, "1001") == false {
		t.Error("1001 wrongly classified as false, expected true")
	}

	// Do the quick check
	f := func(i uint16) bool {
		s := intToBinaryString(i)
		return direct(s) == dfa.Recognize(1, transition, accept, s)
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestDivisibleByN(t *testing.T) {
	// Directly compute divisibility
	direct := func(s uint16, n uint8) bool {
		return s%uint16(n) == 0
	}

	// A DFA that accepts a string of a binary number that is divisible by n.
	makeTransition := func(n int) func(int, byte) int {
		return func(q int, a byte) int {
			// Being in state q for this DFA means that the number we've seen so
			// far has the property i%n == q.
			switch a {
			case '0':
				// Adding a 0 to the end of the string effectively multiplies
				// everything before it by 2 (left shift 1). Then the new number
				// is divisible by 2 and q, or 2*k. Then we take it modulo n
				// since we're ultimately trying to divide by n.
				return (2 * q) % n
			case '1':
				// Similarly, adding a 1 is like multiplying by two (left shift
				// 1) and then adding 1. The transition then follows.
				return (2*q + 1) % n
			default:
				return 0
			}
		}
	}
	accept := map[int]bool{0: true}

	// Func for the quick check
	f := func(i uint16, n uint8) bool {
		if n == 0 {
			// If n is 0, the function is not defined - skip it
			return true
		}

		s := intToBinaryString(i)
		return direct(i, n) == dfa.Recognize(0, makeTransition(int(n)), accept, s)
	}

	// Check for low values of n
	for n := uint8(1); n <= 5; n++ {
		for i := uint16(0); i <= 10; i++ {
			if !f(i, n) {
				t.Error("Failed on input i =", i, "and n =", n)
			}
		}
	}

	// Quick check
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}
