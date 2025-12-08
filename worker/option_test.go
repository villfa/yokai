package worker_test

import (
	"context"
	"testing"

	"github.com/ankorstore/yokai/generate/uuid"
	"github.com/ankorstore/yokai/worker"
	"github.com/ankorstore/yokai/worker/testdata/workers"
	"github.com/stretchr/testify/assert"
)

func TestWorkerPoolOptionsWithGenerator(t *testing.T) {
	t.Parallel()

	generator := uuid.NewDefaultUuidGenerator()

	opt := worker.DefaultWorkerPoolOptions()
	worker.WithGenerator(generator)(&opt)

	assert.Equal(t, generator, opt.Generator)
}

func TestWorkerPoolOptionsWithMetrics(t *testing.T) {
	t.Parallel()

	metrics := worker.NewWorkerMetrics("foo", "bar")

	opt := worker.DefaultWorkerPoolOptions()
	worker.WithMetrics(metrics)(&opt)

	assert.Equal(t, metrics, opt.Metrics)
}

func TestWorkerPoolOptionsWithWorker(t *testing.T) {
	t.Parallel()

	classicWorker := workers.NewClassicWorker()

	registration := worker.NewWorkerRegistration(classicWorker)

	opt := worker.DefaultWorkerPoolOptions()
	worker.WithWorker(classicWorker)(&opt)

	assert.Equal(t, registration, opt.Registrations[classicWorker.Name()])
}

func TestWorkerPoolOptionsWithGlobalDeferredStartThresholds(t *testing.T) {
	t.Parallel()

	opt := worker.DefaultWorkerPoolOptions()
	worker.WithGlobalDeferredStartThreshold(1.5)(&opt)

	assert.Equal(t, 1.5, opt.GlobalDeferredStartThreshold)
}

func TestWorkerPoolOptionsWithGlobalMaxExecutionsAttempts(t *testing.T) {
	t.Parallel()

	opt := worker.DefaultWorkerPoolOptions()
	worker.WithGlobalMaxExecutionsAttempts(2)(&opt)

	assert.Equal(t, 2, opt.GlobalMaxExecutionsAttempts)
}

func TestWorkerExecutionOptionsWithDeferredStartThreshold(t *testing.T) {
	t.Parallel()

	opt := worker.DefaultWorkerExecutionOptions()
	worker.WithDeferredStartThreshold(1.5)(&opt)

	assert.Equal(t, 1.5, opt.DeferredStartThreshold)
}

func TestWorkerExecutionOptionsWithMaxExecutionsAttempts(t *testing.T) {
	t.Parallel()

	opt := worker.DefaultWorkerExecutionOptions()
	worker.WithMaxExecutionsAttempts(2)(&opt)

	assert.Equal(t, 2, opt.MaxExecutionsAttempts)
}

func TestWorkerExecutionOptionsWithMiddleware(t *testing.T) {
	t.Parallel()

	// Test adding a middleware to empty options
	opt1 := worker.DefaultWorkerExecutionOptions()
	worker.WithMiddleware(func(next worker.WorkerMiddlewareFunc) worker.WorkerMiddlewareFunc {
		return func(ctx context.Context) error {
			return next(ctx)
		}
	})(&opt1)

	assert.Len(t, opt1.Middlewares, 1)

	// Test adding a middleware to options with existing middlewares
	opt2 := worker.DefaultWorkerExecutionOptions()

	// Add first middleware directly
	opt2.Middlewares = append(opt2.Middlewares, func(next worker.WorkerMiddlewareFunc) worker.WorkerMiddlewareFunc {
		return func(ctx context.Context) error {
			return next(ctx)
		}
	})

	// Add second middleware using WithMiddleware
	worker.WithMiddleware(func(next worker.WorkerMiddlewareFunc) worker.WorkerMiddlewareFunc {
		return func(ctx context.Context) error {
			return next(ctx)
		}
	})(&opt2)

	assert.Len(t, opt2.Middlewares, 2)

	// Test with nil middlewares slice
	opt3 := worker.ExecutionOptions{
		DeferredStartThreshold: 0,
		MaxExecutionsAttempts:  1,
		Middlewares:            nil,
	}

	worker.WithMiddleware(func(next worker.WorkerMiddlewareFunc) worker.WorkerMiddlewareFunc {
		return func(ctx context.Context) error {
			return next(ctx)
		}
	})(&opt3)

	assert.Len(t, opt3.Middlewares, 1)
}
