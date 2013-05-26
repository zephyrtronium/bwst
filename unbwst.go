package bwst

import (
	"bytes"
	"math/big" // bitset
	"sort"     // TODO: don't use sort to unbwst; interfaces are suboptimal
)

func UnBWST(b []byte) []byte {
	sorted := make([]byte, len(b))
	copy(sorted, b)
	sort.Sort(bytesorter(sorted))
	used := new(big.Int)
	used.SetBit(used, len(b), 1) // reserve capacity
	links := make([]int, len(b))
	// TODO: use O(lg(N)) search in sorted instead of O(N) search in b
	for i, c := range sorted {
		// find the first unused index in b of c
		for j, c2 := range b {
			if c == c2 && used.Bit(j) == 0 {
				links[i] = j
				used.SetBit(used, j, 1)
				break
			}
		}
	}
	// we need to know once again whether each byte is used, so instead of
	// resetting the bitset or using more memory, we can just ask whether it's
	// unused
	unused := used
	words := multibytesorter{}
	for i := range sorted {
		if unused.Bit(i) == 1 {
			word := []byte{}
			x := i
			for unused.Bit(x) == 1 {
				word = append(word, sorted[x])
				unused.SetBit(unused, x, 0)
				x = links[x]
			}
			words = append(words, nil)
			copy(words[1:], words)
			words[0] = word
		}
	}
	if !sort.IsSorted(words) {
		sort.Sort(words)
	}
	x := len(b)
	s := make([]byte, len(b))
	for _, word := range words {
		x -= len(word)
		copy(s[x:], word)
	}
	return s
}

type bytesorter []byte

func (b bytesorter) Len() int           { return len(b) }
func (b bytesorter) Less(i, j int) bool { return b[i] < b[j] }
func (b bytesorter) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }

type multibytesorter [][]byte

func (b multibytesorter) Len() int           { return len(b) }
func (b multibytesorter) Less(i, j int) bool { return bytes.Compare(b[i], b[j]) < 0 }
func (b multibytesorter) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
