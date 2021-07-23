package pools

import (
	"context"
	"math"

	"golang.org/x/sync/semaphore"
)

type TaskPool struct {
	max int
	sem *semaphore.Weighted
}

type Task func() error

func NewTaskPool(max int) *TaskPool {
	size := math.MaxInt32 // 2^32 - 1
	if max > 0 {
		size = max
	}
	return &TaskPool{
		max: size,
		sem: semaphore.NewWeighted(int64(max)),
	}
}

// Run will block until there is available capacity and then execute the given task. Cancelling the context will stop the task from being started.
func (p *TaskPool) Run(ctx context.Context, task Task) <-chan error {
	errc := make(chan error, 1)

	err := p.sem.Acquire(ctx, 1)
	if err != nil {
		errc <- err
		close(errc)
		return errc
	}

	go func() {
		defer p.sem.Release(1)
		defer close(errc)

		err = task()
		if err != nil {
			errc <- err
		}
	}()

	return errc
}

// Wait until all tasks have finished processing.
func (p *TaskPool) Wait() error {
	// acquire all available slots in semaphore
	for i := 0; i < p.max; i++ {
		err := p.sem.Acquire(context.Background(), 1)
		if err != nil {
			return err
		}
	}

	// all tasks have completed; release the semaphore
	p.sem.Release(int64(p.max))

	return nil
}
