package bwst

import "sort"

//TODO: improve sorting
//TODO: e.g. have goroutines sort the word rotations individually in a place
//TODO: which the controller can access, so it can just insert the last-used
//TODO: word into its new location rather than resorting everything

func CyclicLess(a, b []byte, i, j int) bool {
	repeated1, repeated2 := false, false
	x, y := i, j
	for !(repeated1 || repeated2) {
		if a[x] < b[y] {
			return true
		}
		x++
		y++
		if x == len(a) {
			x = 0
		}
		if x == i {
			repeated1 = true
		}
		if y == len(b) {
			y = 0
		}
		if y == j {
			repeated2 = true
		}
	}
	return false
}

type coord struct {
	rotation, word int
}

func minRotations(word []byte, wordN, maxControl int, control chan int, results chan coord) {
	indices := make([]int, len(word))
	for i := range indices {
		indices[i] = i
	}
	indices = sort.Sort(OOPSort{indices, word}).indices // indices of sorted rotations
	for {
		if n, ok := <-control; !ok {
			return
		} else if n < maxControl {
			control <- n + 1
		}
		for i, v := range indices {
			if v >= 0 {
				results <- coord{v, wordN}
				indices[i] = -1
				break
			}
		}
	}
	panic("unreachable")
}

type OOPSort struct {
	indices []int
	word    []byte
}

func (b OOPSort) Len() int {
	return len(b.indices)
}

func (b OOPSort) Swap(i, j int) {
	b.indices[i], b.indices[j] = b.indices[j], b.indices[i]
}

func (b OOPSort) Less(i, j int) bool {
	return CyclicLess(b.word, b.word, b.indices[i], b.indices[j])
}
