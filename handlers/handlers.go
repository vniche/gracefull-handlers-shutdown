package handlers

import (
	"context"
	"sync"

	"github.com/google/uuid"
)

var (
	waitGroup       sync.WaitGroup
	managedHandlers = make(map[string]*Handler)
)

type Handler struct {
	id       string
	quitChan chan struct{}
}

func (handler *Handler) Done() {
	waitGroup.Done()
}

// NewHandler is a constructor for new service handler
func New(ctx context.Context, handlerFunc func(ctx context.Context, shutdown chan struct{}, handler *Handler)) {
	handler := &Handler{
		id:       uuid.New().String(),
		quitChan: make(chan struct{}),
	}
	managedHandlers[handler.id] = handler
	waitGroup.Add(1)
	go handlerFunc(ctx, handler.quitChan, handler)
}

// GracefullyShutdown tries to gracefull shutdown all registered handlers
func GracefullyShutdown() {
	for _, current := range managedHandlers {
		current.quitChan <- struct{}{}
	}

	// wait for all handlers to handle quit signals
	waitGroup.Wait()

	// iterate through handlers and delete their resources
	for ID, current := range managedHandlers {
		if current.quitChan == nil {
			// channel already closed
			return
		}

		// closes handler quit channel
		close(current.quitChan)

		// deletes handler from managed handlers map
		delete(managedHandlers, ID)
	}
}
