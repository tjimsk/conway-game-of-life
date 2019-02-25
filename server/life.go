package life

import (
	"math/rand"
	"time"
)

func init() {
	// for random color generation
	rand.Seed(time.Now().UTC().UnixNano())
}
