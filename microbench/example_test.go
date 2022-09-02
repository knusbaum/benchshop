package main

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func makeRandInts(n int) []int {
	ints := make([]int, 0, n)
	for i := 0; i < n; i++ {
		ints = append(ints, rand.Intn(100))
	}
	return ints
}

func BenchmarkMergesort(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	for _, tt := range []struct {
		in []int
	}{
		{in: makeRandInts(100)},
		{in: makeRandInts(1000)},
		{in: makeRandInts(10000)},
	} {
		data := make([]int, len(tt.in))
		b.Run(fmt.Sprintf("%d", len(tt.in)), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				copy(data, tt.in)
				b.StartTimer()
				mergesort(data)
			}
		})

	}

}

// go test -v -run XXX -bench BenchmarkMergesort -benchmem -cpuprofile cpu.pprof -memprofile mem.pprof -count 10 -benchtime 5s | tee third.txt
