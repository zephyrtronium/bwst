package bwst

import "sort"

func BWST(s []byte) []byte {
	words := factorize(s)
	locs := locate(s, words)
	b := make([]byte, 0, len(s))
	for _, charLocs := range locs {
		sortrots(s, words, charLocs)
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

// compute the Lyndon factorization of s (Duval's algo).
// includes both endpoints.
func factorize(s []byte) (bounds []int) {
	k, m := 0, 1
	bounds = []int{0}
	for m < len(s) {
		a, b := s[k], s[m]
		switch {
		case a < b:
			k = bounds[len(bounds)-1]
			m++
		case a > b:
			bounds = append(bounds, bounds[len(bounds)-1]+m-k)
			k, m = m, m+1
		default:
			k++
			m++
		}
	}
	return append(bounds, len(s))
}

// each instance of a character is considered to be at the beginning of a
// rotation of its word, so the locations can be sorted. because each char is
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

func (l locsorter) Less(i, j int) bool {
	loc1, loc2 := l.locs[i], l.locs[j]
	w1 := l.s[l.words[loc1.word]:l.words[loc1.word+1]]
	w2 := l.s[l.words[loc2.word]:l.words[loc2.word+1]]
	x, y := loc1.idx, loc2.idx
	for i := 0; i < lcm(len(w1), len(w2)); i++ {
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

func lcm(a, b int) int {
	a0, b0 := a, b
	for {
		switch {
		case a < b:
			a += a0
		case a > b:
			b += b0
		default:
			return a
		}
	}
}
