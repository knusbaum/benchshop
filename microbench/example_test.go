package main

import (
	"math/rand"
	"time"
	"testing"
)

func BenchmarkQuicksort(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	var ints []int
	for i := 0; i < 1000; i++ {
		ints = append(ints, rand.Intn(100))
	}
	data := make([]int, len(ints))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		copy(data, ints)
		b.StartTimer()
		mergesort(data)
		b.StopTimer()
	}
}

// go test -v -run XXX -bench BenchmarkQuicksort -benchmem -cpuprofile cpu.pprof -memprofile mem.pprof -count 10 -benchtime 5s | tee first.txt
