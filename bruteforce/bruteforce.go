package bruteforce

func CombinationsBruteForce(alphabet string, n int) []string {
	if n <= 0 {
		return nil
	}

	// Copy alphabet into initial product set -- a set of
	// one character sets
	prod := make([]string, len(alphabet))
	for i, char := range alphabet {
		prod[i] = string(char)
	}

	for i := 1; i < n; i++ {
		// The bigger product should be the size of the alphabet times the size of
		// the n-1 size product
		next := make([]string, 0, len(alphabet)*len(prod))

		// Add each char to each word and add it to the new set
		for _, word := range prod {
			for _, char := range alphabet {
				next = append(next, word+string(char))
			}
		}

		prod = next
	}

	return prod
}
