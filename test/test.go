package test

import (
	"math/rand"
	"os"
	"testing"
	"time"
)

func init() {
	// To avoid seed random number generator and seed it once
	// we do it in init function
	rand.Seed(time.Now().UnixNano())
}

// shouldSkip check environment variable CONDUIT_TEST_API and if it's not set it
// will skip the test identified by state t.
// This is needed to avoid expensive system tests by default without having
// build tags.
func shouldSkip(t *testing.T) {
	_, ok := os.LookupEnv("CONDUIT_TEST_API")
	if !ok {
		t.Skip("CONDUIT_TEST_API not set, skipping system tests")
	}
}

// randString generates random string of length n
func randString(n int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, n)

	for i := 0; i < len(b); i++ {
		b[i] = chars[int(rand.Int63())%len(chars)]
	}

	return string(b)
}
