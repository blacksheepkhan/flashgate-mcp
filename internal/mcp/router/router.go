package router

import (
	"sync"

	"github.com/thomasweidner/flashgate-mcp/internal/mcp/handlers"
	"github.com/thomasweidner/flashgate-mcp/internal/protocol"
)

// Router dispatches method calls to registered handlers.
type Router struct {
	handlers map[string]handlers.Handler
	mu       sync.RWMutex
}

// New creates a new router.
func New() *Router {
	return &Router{
		handlers: make(map[string]handlers.Handler),
	}
}

// Register registers a handler.
func (r *Router) Register(handler handlers.Handler) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.handlers[handler.Method()] = handler
}

// Dispatch dispatches a method call.
func (r *Router) Dispatch(method string, ctx handlers.Context, params []byte) (any, *protocol.Error) {
	r.mu.RLock()
	handler, ok := r.handlers[method]
	r.mu.RUnlock()

	if !ok {
		return nil, &protocol.Error{
			Code:    protocol.ErrMethodNotFound,
			Message: "method not found",
		}
	}

	return handler.Handle(ctx, params)
}
