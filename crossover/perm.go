package crossover

import (
	"math/rand"
)

// OrderBasedCrossover takes a slice of a and fills the gaps with values from b
// while preserving order. 12345678 + 26371485 -> **3456** + 2**71*8* -> 27345618
func OrderBasedCrossover(a, b []int) []int {
	if len(a) != len(b) {
		panic("expected inputs to have same length")
	}
	c := make([]int, len(a))
	start, end := rand.Intn(len(c)), rand.Intn(len(c))
	if start > end {
		start, end = end, start
	}
	// take every value between start and end from a
	taken := map[int]bool{}
	for i := range c {
		if start <= i && i < end {
			c[i] = a[i]
			taken[a[i]] = true
		}
	}
	// return index of next untaken value in b
	nextFromB := func() int {
		for bindex := 0; bindex < len(b); bindex++ {
			_, exists := taken[b[bindex]]
			if !exists {
				return bindex
			}
		}
		return -1
	}
	// fill gaps in c with untaken values from b
	for i := range c {
		if i < start || end <= i {
			nextBIndex := nextFromB()
			taken[b[nextBIndex]] = true
			c[i] = b[nextBIndex]
		}
	}
	return c
}