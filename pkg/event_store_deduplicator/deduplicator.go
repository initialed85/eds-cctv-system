package event_store_deduplicator

import (
	"github.com/initialed85/eds-cctv-system/internal/event_store"
)

type Deduplicator struct {
	store event_store.Store
}

func New(path string) Deduplicator {
	return Deduplicator{
		store: event_store.NewStore(path),
	}
}

func (d *Deduplicator) Deduplicate() error {
	err := d.store.Read()
	if err != nil {
		return err
	}

	return d.store.Write()
}
