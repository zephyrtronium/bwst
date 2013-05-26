package bwst

import (
	"sort"
	"sync"
)

// Compute the Burrows-Wheeler-Scott transform of s. This is done
// out-of-place.
func BWST(s []byte) []byte {
	words := factorize(s)
	// Sorting all rotations of all Lyndon words and then choosing the last
	// character of each is the same as choosing the character to the left of
	// each character in its Lyndon word in sorted order. Therefore, we find
	// all locations of each character, sort them all by their rotations, and
	// proceed therein.
	locs := locate(s, words)
	b := make([]byte, 0, len(s))
	var wg sync.WaitGroup
	for _, charLocs := range locs {
		wg.Add(1)
		go func(charLocs []loc) { defer wg.Done(); sortrots(s, words, charLocs) }(charLocs)
	}
	wg.Wait()
	for _, charLocs := range locs {
		for _, l := range charLocs {
			word := s[words[l.word]:words[l.word+1]]
			i := l.idx - 1
			if i < 0 {
				i = len(word) - 1
			}
			b = append(b, word[i])
		}
	}
	return b
}

// Better than actually storing all rotations of all words. Probably.
type loc struct {
	word, idx int
}

func locate(s []byte, words []int) (locs [256][]loc) {
	w := 0
	for i, c := range s {
		if i >= words[w+1] {
			w++
		}
		locs[int(c)] = append(locs[int(c)], loc{w, i - words[w]})
	}
	return locs
}

// Compute the Lyndon factorization of s (Duval's algo).
// Includes both endpoints.
func factorize(s []byte) (bounds []int) {
	k, m := 0, 1
	bounds = []int{0} // should this be [][]byte? []int should be less mem?
	for m < len(s) {
		a, b := s[k], s[m]
		switch {
		case a < b:
			k = bounds[len(bounds)-1]
			m++
		case a > b:
			// Any time char a is less than char b where b is past a, they are
			// in different Lyndon words because rotating a..b right once is
			// a lesser string.
			bounds = append(bounds, bounds[len(bounds)-1]+m-k)
			k, m = m, m+1
		default:
			k++
			m++
		}
	}
	return append(bounds, len(s))
}

// Each instance of a character is considered to be at the beginning of a
// rotation of its word, so the locations can be sorted. Because each char is
// in order already, we only need to sort the occurrences of each char
// separately to sort the entire thing.

func sortrots(s []byte, words []int, locs []loc) {
	l := locsorter{locs, s, words}
	sort.Sort(l)
}

type locsorter struct {
	locs  []loc
	s     []byte
	words []int
}

func (l locsorter) Len() int      { return len(l.locs) }
func (l locsorter) Swap(i, j int) { l.locs[i], l.locs[j] = l.locs[j], l.locs[i] }

// Cyclic order - AXYA < AXY here because AXYAAXYA < AXYAXY
func (l locsorter) Less(i, j int) bool {
	loc1, loc2 := l.locs[i], l.locs[j]
	// get the actual sequences
	w1 := l.s[l.words[loc1.word]:l.words[loc1.word+1]]
	w2 := l.s[l.words[loc2.word]:l.words[loc2.word+1]]
	x, y := loc1.idx, loc2.idx
	n := lcm(len(w1), len(w2))
	for i := 0; i < n; i++ {
		if a, b := w1[x], w2[y]; a < b {
			return true
		} else if a > b {
			return false
		}
		x++
		if x >= len(w1) {
			x = 0
		}
		y++
		if y >= len(w2) {
			y = 0
		}
	}
	// words are equal
	return false
}

func gcd(m, n int) int {
	var tmp int
	for m != 0 {
		tmp = m
		m = n % m
		n = tmp
	}
	return n
}

func lcm(m, n int) int {
	return m / gcd(m, n) * n
}
