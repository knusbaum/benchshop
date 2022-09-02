package main

import (
	"fmt"
	"math/rand"
	"time"
)

func mergesort(s []int) []int {
	if len(s) == 1 {
		return s
	}

	l := len(s) / 2
	s1 := mergesort(s[:l])
	s2 := mergesort(s[l:])
	var ret []int
	for {
		if s1[0] < s2[0] {
			ret = append(ret, s1[0])
			if len(s1) == 1 {
				ret = append(ret, s2...)
				break
			}
			s1 = s1[1:]
		} else {
			ret = append(ret, s2[0])
			if len(s2) == 1 {
				ret = append(ret, s1...)
				break
			}
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
