package workers

import (
	"context"
	"sync"

	"github.com/ankorstore/yokai/worker"
)

// MiddlewareTestWorker is a worker that records middleware execution
type MiddlewareTestWorker struct {
	mu                sync.Mutex
	middlewaresCalled []string
}

func NewMiddlewareTestWorker() *MiddlewareTestWorker {
	return &MiddlewareTestWorker{
		middlewaresCalled: []string{},
	}
}

func (w *MiddlewareTestWorker) Name() string {
	return "MiddlewareTestWorker"
}

func (w *MiddlewareTestWorker) Run(ctx context.Context) error {
	worker.CtxLogger(ctx).Info().Msg("running middleware test worker")
	return nil
}

func (w *MiddlewareTestWorker) RecordMiddlewareCall(name string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.middlewaresCalled = append(w.middlewaresCalled, name)
}

func (w *MiddlewareTestWorker) GetMiddlewareCalls() []string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return append([]string{}, w.middlewaresCalled...)
}

// CreateTestMiddleware creates a middleware that records its execution
func CreateTestMiddleware(name string, testWorker *MiddlewareTestWorker) worker.WorkerMiddleware {
	return func(next func(ctx context.Context) error) func(ctx context.Context) error {
		return func(ctx context.Context) error {
			// Record that this middleware was called
			testWorker.RecordMiddlewareCall(name)
			
			// Call the next middleware or the worker's Run method
			return next(ctx)
		}
	}
}

// CreateTestMiddlewareWithError creates a middleware that returns an error
func CreateTestMiddlewareWithError(name string, testWorker *MiddlewareTestWorker, err error) worker.WorkerMiddleware {
	return func(next func(ctx context.Context) error) func(ctx context.Context) error {
		return func(ctx context.Context) error {
			// Record that this middleware was called
			testWorker.RecordMiddlewareCall(name)
			
			// Return the error without calling next
			return err
		}
	}
}