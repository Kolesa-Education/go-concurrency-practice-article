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
