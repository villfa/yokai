package worker_test

import (
	"context"
	"testing"

	"github.com/ankorstore/yokai/worker"
	"github.com/ankorstore/yokai/worker/testdata/workers"
	"github.com/stretchr/testify/assert"
)

func TestNewWorkerRegistration(t *testing.T) {
	t.Parallel()

	classicWorker := workers.NewClassicWorker()
	options := []worker.WorkerExecutionOption(nil)

	resolvedWorker := worker.NewWorkerRegistration(classicWorker, options...)

	assert.IsType(t, &worker.WorkerRegistration{}, resolvedWorker)
	assert.Equal(t, classicWorker, resolvedWorker.Worker())
	assert.Equal(t, options, resolvedWorker.Options())
}

func TestWorkerRegistrationMiddlewares(t *testing.T) {
	t.Parallel()

	classicWorker := workers.NewClassicWorker()
	registration := worker.NewWorkerRegistration(classicWorker)

	// Initially, middlewares should be empty
	assert.Empty(t, registration.Middlewares())

	// Add a middleware
	registration.AddMiddleware(func(next worker.WorkerMiddlewareFunc) worker.WorkerMiddlewareFunc {
		return func(ctx context.Context) error {
			return next(ctx)
		}
	})

	// Verify middleware was added
	assert.Len(t, registration.Middlewares(), 1)

	// Add a second middleware
	registration.AddMiddleware(func(next worker.WorkerMiddlewareFunc) worker.WorkerMiddlewareFunc {
		return func(ctx context.Context) error {
			return next(ctx)
		}
	})

	// Verify both middlewares are in the list
	assert.Len(t, registration.Middlewares(), 2)

	// Verify method chaining works
	registration.AddMiddleware(func(next worker.WorkerMiddlewareFunc) worker.WorkerMiddlewareFunc {
		return func(ctx context.Context) error {
			return next(ctx)
		}
	}).AddMiddleware(func(next worker.WorkerMiddlewareFunc) worker.WorkerMiddlewareFunc {
		return func(ctx context.Context) error {
			return next(ctx)
		}
	})

	// Verify middlewares were added via chaining
	assert.Len(t, registration.Middlewares(), 4)
}
