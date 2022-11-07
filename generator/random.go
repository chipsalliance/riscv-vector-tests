package generator

import "math/rand"

// genRandomData generates random bytes by n.
// In order to ensure the determinism of the generated test,
// this function generates pseudo-random data, and for the same n,
// the generated data is always the same.
func genRandomData(n int64) []byte {
	rand.Seed(n)
	data := make([]byte, n)
	rand.Read(data)
	return data
}
