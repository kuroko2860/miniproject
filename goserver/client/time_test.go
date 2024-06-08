package main

import (
	"math/rand"
	"testing"
)

func BenchmarkUnixTimestamp(b *testing.B) {
	types := []string{"car", "bike", "truck", "bus"}
	for i := 0; i < b.N; i++ {
		_ = types[rand.Intn(len(types))]
	}
}
