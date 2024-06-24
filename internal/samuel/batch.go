package samuel

import (
	"encoding/json"
	"math"
)

type Batchable interface {
	json.Marshaler
}

type Batch[T Batchable] struct {
	number int
	count  int
	size   int
	total  int64
	name   string
	items  []T
}

func newBatch[T Batchable](number int, size int, total int64, name string, items ...T) *Batch[T] {
	return &Batch[T]{
		number: number,
		count:  int(math.Ceil(float64(total) / float64(size))),
		size:   size,
		total:  total,
		name:   name,
		items:  items,
	}
}

func (batch *Batch[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"number":   batch.number,
		"count":    batch.count,
		"size":     batch.size,
		"total":    batch.total,
		batch.name: batch.items,
	})
}
