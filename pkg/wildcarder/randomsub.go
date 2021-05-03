package wildcarder

import (
	"math/rand"
	"time"
)

func newRandomSubdomains(count int) []string {
	const letters = "abcdefghijklmnopqrstuvwxyz1234567890"

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var subs []string

	for i := 0; i < count; i++ {
		b := make([]byte, 16)

		for i := range b {
			b[i] = letters[rng.Intn(len(letters))]
		}

		subs = append(subs, string(b))
	}

	return subs
}
