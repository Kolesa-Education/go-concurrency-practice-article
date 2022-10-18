package main

import (
	"crypto/sha256"
	"encoding/hex"
	"go-concurrency-example/bruteforce"
	"log"
	"math/rand"
	"time"
)

const MaxPinSize = 8
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

func searchForCollision(hash string, pinSize int, collisionChan chan string) {
	log.Printf("Iterating %d-sized pins", pinSize)
	combinations := bruteforce.CombinationsBruteForce(allowedPinCharacters, pinSize)
	processPart := func(ccs []string, cc chan string) {
		for _, comb := range combinations {
			bfHash := hexSha256(comb)
			//log.Printf("computing hash for %s:%s", ccs, bfHash)
			if bfHash == hash {
				cc <- comb
			}
		}
	}

	go processPart(combinations[0:len(combinations)/2], collisionChan)
	go processPart(combinations[len(combinations)/2:], collisionChan)
}

func findCollision(hash string, maxPinSize int) string {
	var collisionChan = make(chan string)
	for i := 0; i <= maxPinSize; i++ {
		go searchForCollision(hash, i, collisionChan)
	}
	select {
	case c := <-collisionChan:
		return c
	}
}

func main() {
	size := 8
	pin := randomPinCode(size)
	hash := hexSha256(pin)
	log.Printf("Calculated hash: %s\n", hash)

	start := time.Now()
	collision := findCollision(hash, MaxPinSize)
	if collision == "" {
		log.Printf("couldn't find a collision")
	} else {
		log.Printf("found collision! %s produces hash %s\n", collision, hash)
	}
	end := time.Now().Sub(start)
	log.Printf("Finished in %d ns / %d ms / %ds", end, end/time.Millisecond, end/time.Second)
}
