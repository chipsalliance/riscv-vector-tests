package generator

import "math/rand"

// In order to ensure the determinism of the generated test,
// functions in this file generates pseudo-random data,
// and for the same n, the generated data is always the same.

func genRandomData(n int64) []byte {
	rand.Seed(n)
	data := make([]byte, n)
	rand.Read(data)
	return data
}

func genShuffledSlice(n int) []int {
	s := make([]int, n)
	for i := 0; i < n; i++ {
		s[i] = i
	}
	rand.Seed(int64(n))
	for i := range s {
		j := rand.Intn(i + 1)
		s[i], s[j] = s[j], s[i]
	}
	return s
}
