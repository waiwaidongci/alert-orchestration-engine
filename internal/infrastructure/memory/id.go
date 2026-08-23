package memory

import (
	"fmt"
	"sync/atomic"
	"time"
)

type IDs struct{ n uint64 }

func (i *IDs) NewID(prefix string) string {
	return fmt.Sprintf("%s_%d_%d", prefix, time.Now().UnixNano(), atomic.AddUint64(&i.n, 1))
}
