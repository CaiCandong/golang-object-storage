package main

import (
	"fmt"
	"math/rand"

	"github.com/klauspost/reedsolomon"
)

func fillRandom(p []byte) {
	for i := 0; i < len(p); i += 7 {
		val := rand.Int63()
		for j := 0; i+j < len(p) && j < 7; j++ {
			p[i+j] = byte(val)
			val >>= 8
		}
	}
}

func main() {
	// Create some sample data
	var data = make([]byte, 250000)
	fillRandom(data)

	// Create 5 data slices of 50000 elements each
	enc, _ := reedsolomon.New(5, 3)
	shards, _ := enc.Split(data)
	fmt.Println("shards length is", len(shards))
	for _, shard := range shards {
		fmt.Println(len(shard))
	}
	err := enc.Encode(shards)
	if err != nil {
		panic(err)
	}

	// Check that it verifies
	ok, err := enc.Verify(shards)
	if ok && err == nil {
		fmt.Println("encode ok")
	}

	// Split the data set of 50000 elements into two of 25000
	splitA := make([][]byte, 8)
	splitB := make([][]byte, 8)

	// Merge into a 100000 element set
	merged := make([][]byte, 8)

	// Split/merge the shards
	for i := range shards {
		splitA[i] = shards[i][:25000]
		splitB[i] = shards[i][25000:]

		// Concencate it to itself
		merged[i] = append(make([]byte, 0, len(shards[i])*2), shards[i]...)
		merged[i] = append(merged[i], shards[i]...)
	}

	// Each part should still verify as ok.
	ok, err = enc.Verify(splitA)
	if ok && err == nil {
		fmt.Println("splitA ok")
	}

	ok, err = enc.Verify(splitB)
	if ok && err == nil {
		fmt.Println("splitB ok")
	}

	ok, err = enc.Verify(merged)
	if ok && err == nil {
		fmt.Println("merge ok")
	}
}
