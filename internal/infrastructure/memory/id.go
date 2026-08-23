package memory

import (
	"fmt"
	"time"
)

type IDs struct{ n uint64 }

func (i *IDs) NewID(prefix string) string {
	i.n++
	return fmt.Sprintf("%s_%d_%d", prefix, time.Now().UnixNano(), i.n)
}
