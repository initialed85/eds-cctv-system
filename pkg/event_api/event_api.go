package event_api

import (
	"github.com/initialed85/eds-cctv-system/internal/api"
	"github.com/initialed85/eds-cctv-system/internal/event_store"
)

type API struct {
	store event_store.Store
	api   *api.API
}

func New(path string, port int) *API {
	a := API{
		store: event_store.NewStore(path),
		api:   api.New(port),
	}

	return &a
}
