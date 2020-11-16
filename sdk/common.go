package sdk

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var processBuffer map[string]string = map[string]string{}
var mtx sync.Mutex

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

// SetProcessBuffer will set the current process buffer.
func SetProcessBuffer(ctx, buff string) {

	mtx.Lock()
	defer mtx.Unlock()

	processBuffer[ctx] = buff
}

func getProcessBuffer(ctx string) string {
	mtx.Lock()
	defer mtx.Unlock()

	if buff, ok := processBuffer[ctx]; ok {
		delete(processBuffer, ctx)
		return buff
	}

	return ""
}

func getContext() string {
	return fmt.Sprintf("%v", rand.Int63())
}
