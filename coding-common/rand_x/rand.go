package rand_x

import (
	"math/rand"
	"sync"
	"time"
)

var (
	randSource = rand.NewSource(time.Now().UnixNano())
	randMutex  = &sync.Mutex{}
)

// RandomInt [from, to)
func RandomInt(from, to int) int {
	if from > to {
		from, to = to, from
	}
	randMutex.Lock()
	defer randMutex.Unlock()

	r := rand.New(randSource)
	return r.Intn(to-from+1) + from
}
