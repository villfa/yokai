package worker_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/ankorstore/yokai/generate/generatetest/uuid"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/ankorstore/yokai/worker"
	"github.com/stretchr/testify/assert"
)

// MiddlewareTestWorker is a worker that records middleware execution.
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

// CreateTestMiddleware creates a middleware that records its execution.
func CreateTestMiddleware(name string, testWorker *MiddlewareTestWorker) worker.WorkerMiddleware {
	return func(next worker.WorkerMiddlewareFunc) worker.WorkerMiddlewareFunc {
		return func(ctx context.Context) error {
			// Record that this middleware was called
			testWorker.RecordMiddlewareCall(name)

			// Call the next middleware or the worker's Run method
			return next(ctx)
		}
	}
}

// CreateTestMiddlewareWithError creates a middleware that returns an error.
func CreateTestMiddlewareWithError(name string, testWorker *MiddlewareTestWorker, err error) worker.WorkerMiddleware {
	return func(next worker.WorkerMiddlewareFunc) worker.WorkerMiddlewareFunc {
		return func(ctx context.Context) error {
			// Record that this middleware was called
			testWorker.RecordMiddlewareCall(name)

			// Return the error without calling next
			return err
		}
	}
}

func TestExecutionWithMiddleware(t *testing.T) {
	t.Parallel()

	// test logger
	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, _ := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)

	// test generator
	generator := uuid.NewTestUuidGenerator(testExecutionId)

	// Create a test worker that records middleware execution
	testWorker := NewMiddlewareTestWorker()

	// Create test middlewares
	middleware1 := CreateTestMiddleware("middleware1", testWorker)
	middleware2 := CreateTestMiddleware("middleware2", testWorker)

	// Create the pool with the worker and middlewares
	pool, err := worker.NewDefaultWorkerPoolFactory().Create(
		worker.WithWorker(
			testWorker,
			worker.WithMiddleware(middleware1),
			worker.WithMiddleware(middleware2),
		),
		worker.WithGenerator(generator),
	)
	assert.NoError(t, err)

	// Start the pool
	err = pool.Start(logger.WithContext(context.Background()))
	assert.NoError(t, err)

	// Wait for the worker to complete
	time.Sleep(30 * time.Millisecond)

	// Stop the pool
	err = pool.Stop()
	assert.NoError(t, err)

	// Verify that the middlewares were executed
	middlewareCalls := testWorker.GetMiddlewareCalls()
	assert.Len(t, middlewareCalls, 2)

	// Verify that the middlewares were executed in the correct order
	// Middlewares are executed in the order they are registered
	assert.Equal(t, "middleware1", middlewareCalls[0])
	assert.Equal(t, "middleware2", middlewareCalls[1])

	// Verify the worker execution status
	execution, err := pool.Execution(testWorker.Name())
	assert.NoError(t, err)
	assert.Equal(t, worker.Success, execution.Status())
}

func TestExecutionWithMiddlewareError(t *testing.T) {
	t.Parallel()

	// test logger
	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, _ := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)

	// test generator
	generator := uuid.NewTestUuidGenerator(testExecutionId)

	// Create a test worker that records middleware execution
	testWorker := NewMiddlewareTestWorker()

	// Create a middleware that returns an error
	testError := errors.New("middleware error")
	errorMiddleware := CreateTestMiddlewareWithError("error-middleware", testWorker, testError)

	// Create the pool with the worker and error middleware
	pool, err := worker.NewDefaultWorkerPoolFactory().Create(
		worker.WithWorker(
			testWorker,
			worker.WithMiddleware(errorMiddleware),
		),
		worker.WithGenerator(generator),
	)
	assert.NoError(t, err)

	// Start the pool
	err = pool.Start(logger.WithContext(context.Background()))
	assert.NoError(t, err)

	// Wait for the worker to complete
	time.Sleep(30 * time.Millisecond)

	// Stop the pool
	err = pool.Stop()
	assert.NoError(t, err)

	// Verify that the middleware was executed
	middlewareCalls := testWorker.GetMiddlewareCalls()
	assert.Len(t, middlewareCalls, 1)
	assert.Equal(t, "error-middleware", middlewareCalls[0])

	// Verify the worker execution status
	execution, err := pool.Execution(testWorker.Name())
	assert.NoError(t, err)
	assert.Equal(t, worker.Error, execution.Status())

	// Verify that the error was logged
	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":             "error",
		"worker":            testWorker.Name(),
		"workerExecutionID": testExecutionId,
		"error":             "middleware error",
	})
}

func TestExecutionWithMiddlewareChaining(t *testing.T) {
	t.Parallel()

	// test logger
	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, _ := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)

	// test generator
	generator := uuid.NewTestUuidGenerator(testExecutionId)

	// Create a test worker that records middleware execution
	testWorker := NewMiddlewareTestWorker()

	// Create test middlewares
	middleware1 := CreateTestMiddleware("middleware1", testWorker)
	middleware2 := CreateTestMiddleware("middleware2", testWorker)
	middleware3 := CreateTestMiddleware("middleware3", testWorker)

	// Create the pool with the worker and middlewares
	pool, err := worker.NewDefaultWorkerPoolFactory().Create(
		worker.WithWorker(
			testWorker,
			worker.WithMiddleware(middleware1),
			worker.WithMiddleware(middleware2),
			worker.WithMiddleware(middleware3),
		),
		worker.WithGenerator(generator),
	)
	assert.NoError(t, err)

	// Start the pool
	err = pool.Start(logger.WithContext(context.Background()))
	assert.NoError(t, err)

	// Wait for the worker to complete
	time.Sleep(30 * time.Millisecond)

	// Stop the pool
	err = pool.Stop()
	assert.NoError(t, err)

	// Verify that the middlewares were executed
	middlewareCalls := testWorker.GetMiddlewareCalls()
	assert.Len(t, middlewareCalls, 3)

	// Verify that the middlewares were executed in the correct order
	// Middlewares are executed in the order they are registered
	assert.Equal(t, "middleware1", middlewareCalls[0])
	assert.Equal(t, "middleware2", middlewareCalls[1])
	assert.Equal(t, "middleware3", middlewareCalls[2])

	// Verify the worker execution status
	execution, err := pool.Execution(testWorker.Name())
	assert.NoError(t, err)
	assert.Equal(t, worker.Success, execution.Status())
}
