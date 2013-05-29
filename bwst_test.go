package bwst

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"testing"
)

const test_string = `Wherever you go, I follow,
hands held across sidewalks
frozen
or slopes more slippery
called love,
life, and 1685 miles
of "I can't wait
to be there when you fall."`

func TestAbsorption(t *testing.T) {
	s := []byte(test_string)
	s = UnBWST(BWST(s))
	if string(s) != test_string {
		t.Fatal("UnBWST(BWST(s)) failed: expected a sweet poem, got gibberish")
	}
	s = BWST(UnBWST(s))
	if string(s) != test_string {
		t.Fatal("BWST(UnBWST(s)) failed: expected a sweet poem, got gibberish")
	}
}

func TestAbsorptionRandom(t *testing.T) {
	s := makerandombytes(1 << 15)
	c := s
	s = UnBWST(BWST(s))
	if !bytes.Equal(c, s) {
		t.Fatal("UnBWST(BWST(s)) failed: expected randomness, got gibberish")
	}
	s = BWST(UnBWST(s))
	if !bytes.Equal(c, s) {
		t.Fatal("BWST(UnBWST(s)) failed: expected randomness, got gibberish")
	}
}

func BenchmarkBWST(b *testing.B) {
	s := []byte(test_string)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = BWST(s)
	}
}

func BenchmarkBWSTRandom(b *testing.B) {
	s := makerandombytes(1 << 15)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = BWST(s)
	}
}

func BenchmarkUnBWST(b *testing.B) {
	s := []byte(test_string)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = UnBWST(s)
	}
}

func BenchmarkUnBWSTRandom(b *testing.B) {
	s := makerandombytes(1 << 15)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = UnBWST(s)
	}
}

func makerandombytes(n int) (b []byte) {
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	c, _ := aes.NewCipher(iv) // I do insist.
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	s := cipher.NewCTR(c, iv)
	b = make([]byte, n)
	s.XORKeyStream(b, b)
	return b
}
