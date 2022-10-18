package bruteforce

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestCombinationsBruteForce(t *testing.T) {
	t.Run("ABC, 2", func(t *testing.T) {
		alphabet := "ABC"
		combinations := CombinationsBruteForce(alphabet, 2)
		expected := []string{
			"AA", "AB", "AC", "BA", "BB", "BC", "CA", "CB", "CC",
		}
		log.Println(combinations)
		assert.Equal(t, expected, combinations)
	})
	t.Run("ABC, 3", func(t *testing.T) {
		alphabet := "ABC"
		combinations := CombinationsBruteForce(alphabet, 3)
		expected := []string{
			"AAA", "AAB", "AAC", "ABA", "ABB", "ABC", "ACA", "ACB", "ACC",
			"BAA", "BAB", "BAC", "BBA", "BBB", "BBC", "BCA", "BCB", "BCC",
			"CAA", "CAB", "CAC", "CBA", "CBB", "CBC", "CCA", "CCB", "CCC",
		}
		log.Println(combinations)
		assert.Equal(t, expected, combinations)
	})

}
