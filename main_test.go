package main

import (
	"github.com/stretchr/testify/assert"
	"log"
	"math/rand"
	"testing"
)

func Test_randomPinCode(t *testing.T) {
	t.Run("size = 10, seed = 1", func(t *testing.T) {
		var seed int64 = 1
		r := rand.New(rand.NewSource(seed))
		size := 10

		pin := randomPinCodeWithRand(size, *r)
		log.Println(pin)
		assert.Equal(t, 10, len(pin))
		assert.Equal(t, "1779185060", pin)
	})

	t.Run("size = 5, seed = 1", func(t *testing.T) {
		var seed int64 = 1
		r := rand.New(rand.NewSource(seed))
		size := 5

		pin := randomPinCodeWithRand(size, *r)
		log.Println(pin)
		assert.Equal(t, 5, len(pin))
		assert.Equal(t, "17791", pin)
	})
}

func Benchmark_combinations(b *testing.B) {
	b.Run("99999999", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			combinations("99999999")
		}
	})

	b.Run("9999999", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			combinations("9999999")
		}
	})

	b.Run("999999", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			combinations("999999")
		}
	})

	b.Run("99999", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			combinations("99999")
		}
	})

	b.Run("11111111", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			combinations("11111111")
		}
	})

	b.Run("1111111", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			combinations("1111111")
		}
	})

	b.Run("111111", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			combinations("111111")
		}
	})

	b.Run("11111", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			combinations("11111")
		}
	})

	b.Run("55555555", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			combinations("55555555")
		}
	})

	b.Run("5555555", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			combinations("5555555")
		}
	})

	b.Run("555555", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			combinations("555555")
		}
	})

	b.Run("55555", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			combinations("55555")
		}
	})
}
