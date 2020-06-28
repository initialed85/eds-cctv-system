package event_store_deduplicator

import (
	"github.com/initialed85/eds-cctv-system/internal/event_store"
	"log"
)

type Deduplicator struct {
	sourceStore, destinationStore *event_store.Store
}

func New(sourcePath, destinationPath string) Deduplicator {
	return Deduplicator{
		sourceStore:      event_store.NewStore(sourcePath),
		destinationStore: event_store.NewStore(destinationPath),
	}
}

func (d *Deduplicator) Deduplicate() error {
	err := d.sourceStore.Read()
	if err != nil {
		return err
	}

	log.Printf("read %v events from sourceStore", d.sourceStore.Len())

	events := d.sourceStore.GetAll()
	for _, event := range events {
		d.destinationStore.Add(event)
	}

	log.Printf("wrote %v events to destinationStore", d.destinationStore.Len())

	return d.destinationStore.Write()
}
