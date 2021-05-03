package threadpool

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

// ThreadPool is a thread pool object used to execute tasks in parallel.
type ThreadPool struct {
	taskCounter     int64
	taskDoneCounter int64

	taskChan chan Runnable
	doneChan chan bool

	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

// Runnable defines an interface with a Run function to be executed by a thread in a thread pool.
type Runnable interface {
	Run()
}

// NewThreadPool creates a new ThreadPool object and starts the worker threads.
func NewThreadPool(threads int, queueSize int) *ThreadPool {
	p := &ThreadPool{}
	p.taskChan = make(chan Runnable, queueSize)
	p.doneChan = make(chan bool, queueSize)

	p.ctx, p.cancel = context.WithCancel(context.Background())

	p.createPool(threads)
	p.createSentinel()

	return p
}

// Execute adds a task to the task queue to be picked by a worker thread. It will block if the queue is full.
func (p *ThreadPool) Execute(task Runnable) {
	atomic.AddInt64(&p.taskCounter, 1)
	p.taskChan <- task
}

// Done returns true if there are no tasks in flight.
func (p *ThreadPool) Done() bool {
	current := atomic.LoadInt64(&p.taskCounter)
	done := atomic.LoadInt64(&p.taskDoneCounter)
	return current == done
}

// Wait waits for all the tasks in flight to be processed.
func (p *ThreadPool) Wait() {
	for !p.Done() {
		time.Sleep(1 * time.Millisecond)
	}
}

// Close closes the threadpool and frees up the threads.
func (p *ThreadPool) Close() {
	p.Wait()

	p.cancel()
	p.wg.Wait()

	close(p.taskChan)
	close(p.doneChan)
}

// CurrentCount returns the current number of tasks processed.
func (p *ThreadPool) CurrentCount() int {
	done := atomic.LoadInt64(&p.taskDoneCounter)
	return int(done)
}

func (p *ThreadPool) createPool(threads int) {
	for i := 0; i < threads; i++ {
		worker := newWorker(p.ctx, &p.wg, p.taskChan, p.doneChan)
		worker.start()
	}
}

func (p *ThreadPool) createSentinel() {
	go func(ctx context.Context, counter *int64) {
		for {
			select {
			case <-ctx.Done():
				return
			case <-p.doneChan:
				atomic.AddInt64(counter, 1)
			}
		}
	}(p.ctx, &p.taskDoneCounter)
}
