package notification

import "time"

type Batch struct {
	Items     []Record
	CreatedAt time.Time
}

func NewBatch(items []Record, now time.Time) Batch {
	return Batch{Items: append([]Record(nil), items...), CreatedAt: now}
}
func (b Batch) Pending() []Record {
	out := []Record{}
	for _, r := range b.Items {
		if r.Status == Pending {
			out = append(out, r)
		}
	}
	return out
}
func (b Batch) Failed() []Record {
	out := []Record{}
	for _, r := range b.Items {
		if r.Status == Failed {
			out = append(out, r)
		}
	}
	return out
}
func (b Batch) SentCount() int {
	n := 0
	for _, r := range b.Items {
		if r.Status == Sent {
			n++
		}
	}
	return n
}
func (b Batch) FailedCount() int {
	n := 0
	for _, r := range b.Items {
		if r.Status == Failed || r.Status == DeadLetter {
			n++
		}
	}
	return n
}
func (b Batch) Complete() bool {
	return len(b.Items) > 0 && b.SentCount()+b.FailedCount() == len(b.Items)
}
