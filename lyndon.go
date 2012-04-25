package bwst

// Calculate the Lyndon factorization of s. The slices in the returned slice
// map into the memory of s; modifying either modifies the other.
//
// Duval's factorization algorithm: http://en.wikipedia.org/wiki/Lyndon_word
func LyndonBytes(s []byte) (factors [][]byte) {
	k, m := 0, 1
	for m < len(s) {
		a, b := s[k], s[m]
		switch {
		case a < b:
			k = 0
			m++
		case a > b:
			factors = append(factors, s[:m-k])
			s = s[m-k:]
			k, m = 0, 1
		default:
			k++
			m++
		}
	}
	// This conditional converts a single-character run at the end of the
	// string into a single word. This may actually cause the output not to be
	// a true Lyndon factorization, since the last word may not be a true
	// Lyndon word, but for the BWST application, it increases efficiency. To
	// make this a true Lyndon factorization, simply remove this if.
	if m == len(s) && m-k == 1 {
		return append(factors, s)
	}
	for len(s) > 0 {
		factors = append(factors, s[:m-k])
		s = s[m-k:]
	}
	return
}
