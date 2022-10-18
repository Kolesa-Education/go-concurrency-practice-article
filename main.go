package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"go-concurrency-example/bruteforce"
	"log"
	"math/rand"
	"time"
)

const MaxPinSize = 10
const allowedPinCharacters string = "0123456789"

func randomPinCode(size int) string {
	return randomPinCodeWithRand(size, *rand.New(rand.NewSource(time.Now().UnixNano())))
}

func randomPinCodeWithRand(size int, r rand.Rand) string {
	b := make([]byte, size)
	for i := range b {
		b[i] = allowedPinCharacters[r.Intn(len(allowedPinCharacters))]
	}
	return string(b)
}

func hexSha256(input string) string {
	hashedPin := sha256.Sum256([]byte(input))
	hexHashedPin := hex.EncodeToString(hashedPin[:])
	return hexHashedPin
}

func findCollision(hash string, maxPinSize int) (string, error) {
	for i := 0; i < maxPinSize; i++ {
		log.Printf("Iterating %d-sized pins", i)
		combinations := bruteforce.CombinationsBruteForce(allowedPinCharacters, i)
		for _, c := range combinations {
			bfHash := hexSha256(c)
			if bfHash == hash {
				return bfHash, nil
			}
		}
	}
	return "", errors.New("not found")
}

func main() {
	size := 8
	pin := randomPinCode(size)
	hash := hexSha256(pin)
	log.Printf("Calculated hash: %s\n", hash)

	start := time.Now()
	collision, err := findCollision(hash, MaxPinSize)
	if err != nil {
		log.Printf("couldn't find a collision")
	} else {
		log.Printf("found collision! %s produces hash %s\n", collision, hash)
	}
	end := time.Now().Sub(start)
	log.Printf("Finished in %d ns / %d ms / %ds", end, end/time.Millisecond, end/time.Second)
}
