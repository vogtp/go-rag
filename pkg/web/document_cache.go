package web

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type docChace struct {
	mu    sync.RWMutex
	cache sync.Map
}

func newDocCache() docChace {
	return docChace{}
}

func (dc *docChace) add(d *queryDoc) {
	dc.cache.Store(d.UUID, *d)
}

func (dc *docChace) get(id uuid.UUID) (*queryDoc, error) {
	d, ok := dc.cache.LoadAndDelete(id)
	if !ok {
		return nil, fmt.Errorf("cannot find document for %v", id)
	}
	doc,ok:=d.(queryDoc)
	if !ok {
		return nil, fmt.Errorf("document %v is not correct type", id)
	}
	return &doc, nil
}
