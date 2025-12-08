package worker

import "context"

// Worker is the interface to implement to provide workers.
type Worker interface {
	Name() string
	Run(ctx context.Context) error
}

// WorkerMiddleware is a function that wraps a worker's Run method to add behavior.
type WorkerMiddleware func(next WorkerMiddlewareFunc) WorkerMiddlewareFunc

// WorkerMiddlewareFunc is a function type that represents a worker's Run method signature.
// It is used by WorkerMiddleware to wrap the worker execution with additional behavior.
type WorkerMiddlewareFunc func(ctx context.Context) error

// WorkerRegistration is a [Worker] registration, with optional [WorkerExecutionOption] and [WorkerMiddleware].
type WorkerRegistration struct {
	worker      Worker
	options     []WorkerExecutionOption
	middlewares []WorkerMiddleware
}

// NewWorkerRegistration returns a new [WorkerRegistration] for a given [Worker] and an optional list of [WorkerExecutionOption].
func NewWorkerRegistration(worker Worker, options ...WorkerExecutionOption) *WorkerRegistration {
	return &WorkerRegistration{
		worker:      worker,
		options:     options,
		middlewares: []WorkerMiddleware{},
	}
}

// Worker returns the [Worker] of the [WorkerRegistration].
func (r *WorkerRegistration) Worker() Worker {
	return r.worker
}

// Options returns the list of [WorkerExecutionOption] of the [WorkerRegistration].
func (r *WorkerRegistration) Options() []WorkerExecutionOption {
	return r.options
}

// Middlewares returns the list of [WorkerMiddleware] of the [WorkerRegistration].
func (r *WorkerRegistration) Middlewares() []WorkerMiddleware {
	return r.middlewares
}

// AddMiddleware adds a [WorkerMiddleware] to the [WorkerRegistration].
func (r *WorkerRegistration) AddMiddleware(middleware WorkerMiddleware) *WorkerRegistration {
	r.middlewares = append(r.middlewares, middleware)

	return r
}
