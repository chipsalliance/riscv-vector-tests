package generator

import (
	"math/rand"
	"sync"
)

// In order to ensure the determinism of the generated test,
// functions in this file generates pseudo-random data,
// and for the same n, the generated data is always the same.

var mu = sync.Mutex{}

func getRandomInt(n int) int {
	mu.Lock()
	defer mu.Unlock()
	return rand.New(rand.NewSource(int64(n))).Intn(n)
}

func genRandomData(n int64) []byte {
	mu.Lock()
	defer mu.Unlock()
	data := make([]byte, n)
	rand.New(rand.NewSource(n)).Read(data)
	return data
}

func genShuffledSlice(n int) []int {
	mu.Lock()
	defer mu.Unlock()
	s := make([]int, n)
	for i := 0; i < n; i++ {
		s[i] = i
	}
	for i := range s {
		j := rand.New(rand.NewSource(int64(n))).Intn(i + 1)
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func shuffleSlice(s []int, n int64) {
	mu.Lock()
	defer mu.Unlock()
	for i := range s {
		j := rand.New(rand.NewSource(n)).Intn(i + 1)
		s[i], s[j] = s[j], s[i]
	}
}
