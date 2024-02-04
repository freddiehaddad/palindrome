package palindrome

import (
	"log"
	"os"
	"sync"
)

// A char represents a single palindrome character b at position pos
// within a string.
type char struct {
	b   byte
	pos int
}

// init prepares the logger.
func init() {
	log.SetOutput(os.Stderr)
	log.SetFlags(log.LUTC | log.Lshortfile | log.LstdFlags)
}

// drain empties channel c by reading and discarding all remaining elements.
func drain(c chan char) {
	for range c {
	}
}

// toLower returns the lowercase equivalent of b when it's a character between
// A and Z.  Otherwise b is returned unchanged.
func toLower(b byte) byte {
	if b >= 'A' && b <= 'Z' {
		return b + 0x20
	}
	return b
}

// isPalindromeCharacter returns true if b is an alphanumeric character
// based on the ASCII character table.
func isPalindromeCharacter(b byte) bool {
	if b >= '0' && b <= '9' {
		return true
	}

	if b >= 'a' && b <= 'z' {
		return true
	}

	if b >= 'A' && b <= 'Z' {
		return true
	}

	return false
}

// feedCharactersForward reads characters in s from left to right until a stop
// message is received or the end of string is reached. Any palindrome character
// is then written to channel c. The channel c when the function returns.
func feedCharactersForward(wg *sync.WaitGroup, c chan char, stop chan bool, s string) {
	defer wg.Done()

	for i := 0; i < len(s); i++ {
		select {
		case <-stop:
			log.Printf("feedCharactersForward: stop signal")
			return
		default:
			ch := s[i]

			if !isPalindromeCharacter(ch) {
				log.Printf("feedCharactersForward: skipping ch=%c", ch)
				continue
			}

			ch = toLower(ch)
			packet := char{ch, i}
			c <- packet
			log.Printf("feedCharactersForward: sent packet=%v", packet)
		}
	}
}

// feedCharactersReverse reads characters in s from right to left until a stop
// message is received or the end of string is reached. Any palindrome character
// is then written to channel c. The channel c when the function returns.
func feedCharactersReverse(wg *sync.WaitGroup, c chan char, stop chan bool, s string) {
	defer wg.Done()

	for i := len(s) - 1; i >= 0; i-- {
		select {
		case <-stop:
			log.Printf("feedCharactersReverse: stop signal")
			return
		default:
			ch := s[i]

			if !isPalindromeCharacter(ch) {
				log.Printf("feedCharactersReverse: skipping ch=%c", ch)
				continue
			}

			ch = toLower(ch)
			packet := char{ch, i}
			c <- packet
			log.Printf("feedCharactersReverse: sent packet=%v", packet)
		}
	}
}

// compareCharacters reads characters from the left and right channels comparing
// each pair. Comparisons continue until the intersection is reached or two
// different characters are compared. If all characters compared match, true is
// written to the result channel. Otherwise, false is written. The stop and
// result channels are closed when the function returns.
func compareCharacters(wg *sync.WaitGroup, left chan char, right chan char, stop chan bool, result chan bool) {
	defer wg.Done()

	for {
		l := <-left
		r := <-right
		log.Printf("compareCharacters: comparing l=%v r=%v", l, r)

		// Reached the middle of the string
		if l.pos >= r.pos {
			result <- true
			log.Printf("compareCharacters: intersection l=%v r=%v", l, r)
			break
		}

		// Compare characters
		if l.b != r.b {
			result <- false
			log.Printf("compareCharacters: mismatch l=%v r=%v", l, r)
			break
		}
	}
	stop <- true
	stop <- true

	// Drain channels
	go drain(left)
	go drain(right)
}

// isPalindrome examines the string s and determines if the right half is a
// mirror reflection of the left half.  Non-palindrome characters are ignored
// (see isPalindromeCharacter). If s is determined to be a palindrome, true is
// returned, otherwise false.
func isPalindrome(s string) bool {
	if len(s) <= 1 {
		return true
	}

	wg := &sync.WaitGroup{}
	wg.Add(3)

	// result from compareCharacters
	result := make(chan bool, 1)

	// signal to stop comparing for character feeders
	stop := make(chan bool, 2)

	// channels for feeding characters to comparitor
	leftChannel := make(chan char, 1)
	rightChannel := make(chan char, 1)

	defer close(stop)
	defer close(result)
	defer close(leftChannel)
	defer close(rightChannel)

	// go routines for processing the string
	go feedCharactersForward(wg, leftChannel, stop, s)
	go feedCharactersReverse(wg, rightChannel, stop, s)
	go compareCharacters(wg, leftChannel, rightChannel, stop, result)

	wg.Wait()
	r := <-result
	return r
}
