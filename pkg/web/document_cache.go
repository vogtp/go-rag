package web

import (
	"fmt"

	"github.com/google/uuid"
)

type docChace struct {
	cache map[uuid.UUID]queryDoc
}

func newDocCache() docChace {
	return docChace{
		cache: make(map[uuid.UUID]queryDoc),
	}
}

func (dc *docChace) add(d *queryDoc) {
	dc.cache[d.UUID] = *d
}

func (dc *docChace) get(id uuid.UUID) (*queryDoc, error) {
	d, ok := dc.cache[id]
	if !ok {
		return nil, fmt.Errorf("cannot find document for %v", id)
	}
	delete(dc.cache, id)
	return &d, nil
}
