package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Partition chooses a pivot and puts all items < pivot before the pivot and all items > pivot
// after pivot. It then returns pivot.
func partition(s []int, start, end int) int {
	pivot := s[end]

	i := start - 1

	for j := start; j < end; j++ {
		if s[j] <= pivot {
			i += 1
			s[i], s[j] = s[j], s[i]
		}
	}
	i += 1
	s[i], s[end] = s[end], s[i]
	return i
}

func quicksortr(s []int, start, end int) {
	if start >= end || start < 0 {
		return
	}

	p := partition(s, start, end)
	quicksortr(s, start, p-1)
	quicksortr(s, p+1, end)
}

func quicksort(s []int) {
	quicksortr(s, 0, len(s)-1)
}

func mergesort(s []int) []int {
	if len(s) == 1 {
		return s
	}

	l := len(s) / 2
	s1 := mergesort(s[:l])
	s2 := mergesort(s[l:])
	var ret []int
	for {
		if len(s1) == 0 {
			ret = append(ret, s2...)
			break
		} else if len(s2) == 0 {
			ret = append(ret, s1...)
			break
		}
		if s1[0] < s2[0] {
			ret = append(ret, s1[0])
			s1 = s1[1:]
		} else {
			ret = append(ret, s2[0])
			s2 = s2[1:]
		}
	}
	return ret
}

func main() {
	rand.Seed(time.Now().UnixNano())
	var ints []int
	for i := 0; i < 10; i++ {
		ints = append(ints, rand.Intn(100))
	}
	fmt.Printf("Before:\n%v\n", ints)
	ints = mergesort(ints)
	fmt.Printf("After:\n%v\n", ints)
}
