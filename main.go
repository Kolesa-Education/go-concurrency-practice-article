package main

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/Kolesa-Education/go-concurrency-practice-article/bruteforce"
	"log"
	"math/rand"
	"runtime"
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
	for _, comb := range combinations {
		bfHash := hexSha256(comb)
		//log.Printf("computing hash for %s:%s", ccs, bfHash)
		if bfHash == hash {
			collisionChan <- comb
		}
	}
}

func findCollision(hash string, maxPinSize int, maxGoroutines int) string {
	guard := make(chan any, maxGoroutines)
	var collisionChan = make(chan string, maxGoroutines)

	for i := 0; i <= maxPinSize; i++ {
		guard <- struct{}{}
		go func(i int) {
			searchForCollision(hash, i, collisionChan)
			<-guard
		}(i)
	}
	select {
	case c := <-collisionChan:
		return c
	}
}

func measure(f func()) time.Duration {
	start := time.Now()
	f()
	end := time.Now().Sub(start)
	return end
}

func combinations(pin string) {
	hash := hexSha256(pin)
	log.Printf("Calculated hash: %s\n", hash)
	duration := measure(func() {
		collision := findCollision(hash, MaxPinSize, runtime.NumCPU())
		if collision == "" {
			log.Printf("couldn't find a collision")
		} else {
			log.Printf("found collision! %s produces hash %s\n", collision, hash)
		}
	})
	log.Printf("Finished in %d ns / %d ms / %ds", duration, duration/time.Millisecond, duration/time.Second)
}

func main() {
	size := 8
	pin := randomPinCode(size)
	log.Printf("runtime cores accessible %d\n", runtime.NumCPU())
	combinations(pin)

}
