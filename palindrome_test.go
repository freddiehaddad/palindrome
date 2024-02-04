package palindrome

import (
	"sync"
	"testing"
)

func TestToLower(t *testing.T) {
	tests := []struct {
		input    byte
		expected byte
	}{
		{'@', '@'},
		{'[', '['},
		{'`', '`'},
		{'{', '{'},
		{'a', 'a'},
		{'A', 'a'},
		{'b', 'b'},
		{'B', 'b'},
		{'y', 'y'},
		{'Y', 'y'},
		{'z', 'z'},
		{'Z', 'z'},
		{'0', '0'},
		{'9', '9'},
		{'!', '!'},
		{' ', ' '},
	}

	for i, test := range tests {
		result := toLower(test.input)
		if result != test.expected {
			t.Error("Test[", i, "]:", "failed.", "expected", test.expected, "got", result)
		}
	}
}

func TestIsPalindromeCharacter(t *testing.T) {
	tests := []struct {
		input    byte
		expected bool
	}{
		{'@', false},
		{'[', false},
		{'`', false},
		{'{', false},
		{'/', false},
		{':', false},
		{'a', true},
		{'A', true},
		{'z', true},
		{'Z', true},
		{'0', true},
		{'9', true},
		{' ', false},
		{'!', false},
		{',', false},
		{'%', false},
	}

	for i, test := range tests {
		result := isPalindromeCharacter(test.input)
		if result != test.expected {
			t.Error("Test[", i, "]:", "failed.", "expected", test.expected, "got", result)
		}
	}
}

func TestFeedCharactersForward(t *testing.T) {
	tests := []struct {
		input    string
		expected []char
	}{
		{"", []char{}},
		{" ", []char{}},
		{"a", []char{{'a', 0}}},
		{"a ", []char{{'a', 0}}},
		{" a", []char{{'a', 1}}},
		{"a b", []char{{'a', 0}, {'b', 2}}},
		{"a b ", []char{{'a', 0}, {'b', 2}}},
		{" a b", []char{{'a', 1}, {'b', 3}}},
		{" a b ", []char{{'a', 1}, {'b', 3}}},
		{
			"Madam, I'm Adam.",
			[]char{
				{'m', 0},
				{'a', 1},
				{'d', 2},
				{'a', 3},
				{'m', 4},
				{'i', 7},
				{'m', 9},
				{'a', 11},
				{'d', 12},
				{'a', 13},
				{'m', 14},
			},
		},
	}

	for i, test := range tests {
		wg := &sync.WaitGroup{}
		ch := make(chan char)
		stop := make(chan bool)

		wg.Add(1)
		go feedCharactersForward(wg, ch, stop, test.input)

		for j, b := range test.expected {
			c := <-ch

			if b.b != c.b {
				t.Error("Test[", i, ":", j, "]:", "failed.", "expected", string(b.b), "got", string(c.b))
			}

			if b.pos != c.pos {
				t.Error("Test[", i, ":", j, "]:", "failed.", "expected", "pos", b.pos, "got", "pos", c.pos)
			}
		}

		close(ch)
		close(stop)
	}
}

func TestFeedCharactersReverse(t *testing.T) {
	tests := []struct {
		input    string
		expected []char
	}{
		{"", []char{}},
		{" ", []char{}},
		{"a", []char{{'a', 0}}},
		{"a ", []char{{'a', 0}}},
		{" a", []char{{'a', 1}}},
		{"a b", []char{{'b', 2}, {'a', 0}}},
		{"a b ", []char{{'b', 2}, {'a', 0}}},
		{" a b", []char{{'b', 3}, {'a', 1}}},
		{" a b ", []char{{'b', 3}, {'a', 1}}},
		{
			"Madam, I'm Adam.",
			[]char{
				{'m', 14},
				{'a', 13},
				{'d', 12},
				{'a', 11},
				{'m', 9},
				{'i', 7},
				{'m', 4},
				{'a', 3},
				{'d', 2},
				{'a', 1},
				{'m', 0},
			},
		},
	}

	for i, test := range tests {
		wg := &sync.WaitGroup{}
		ch := make(chan char)
		stop := make(chan bool)

		wg.Add(1)
		go feedCharactersReverse(wg, ch, stop, test.input)

		for j, b := range test.expected {
			c := <-ch

			if b.b != c.b {
				t.Error("Test[", i, ":", j, "]:", "failed.", "expected", string(b.b), "got", string(c.b))
			}

			if b.pos != c.pos {
				t.Error("Test[", i, ":", j, "]:", "failed.", "expected", "pos", b.pos, "got", "pos", c.pos)
			}
		}

		close(ch)
		close(stop)
	}
}

func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"", true},
		{"a", true},
		{"aa", true},
		{"ab", false},
		{"aba", true},
		{"abc", false},
		{"aab", false},
		{"abb", false},
		{"abba", true},
		{"abbc", false},
		{"abca", false},
		{"aaba", false},
		{"abaa", false},
		{"a,baa", false},
		{"ab,aa", false},
		{"aba,a", false},
		{"aa,aa", true},
		{"   aa", true},
		{"aa   ", true},
		{"aaabaaacaaa", false},
		{"A man, a plan, a canal: Panama", true},
	}

	for _, test := range tests {
		result := isPalindrome(test.input)
		if result != test.expected {
			t.Error("Test[", test.input, "]:", "failed.", "expected", test.expected, "got", result)
		}
	}
}

func BenchmarkIsPalindrome(b *testing.B) {
	s := "A man, a plan, a canal: Panama"
	for i := 0; i < b.N; i++ {
		isPalindrome(s)
	}
}
