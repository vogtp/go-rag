package web

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type docChace struct {
	mu sync.Mutex
	cache map[uuid.UUID]queryDoc
}

func newDocCache() docChace {
	return docChace{
		cache: make(map[uuid.UUID]queryDoc),
	}
}

func (dc *docChace) add(d *queryDoc) {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	dc.cache[d.UUID] = *d
}

func (dc *docChace) get(id uuid.UUID) (*queryDoc, error) {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	d, ok := dc.cache[id]
	if !ok {
		return nil, fmt.Errorf("cannot find document for %v", id)
	}
	delete(dc.cache, id)
	return &d, nil
}
