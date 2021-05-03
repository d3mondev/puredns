package threadpool

import (
	"context"
)

type worker struct {
	taskChan chan Runnable
	doneChan chan bool
	ctx      context.Context
	wg       signaler
}

type signaler interface {
	Add(delta int)
	Done()
	Wait()
}

func newWorker(ctx context.Context, wg signaler, taskChan chan Runnable, doneChan chan bool) *worker {
	wg.Add(1)

	return &worker{
		ctx:      ctx,
		wg:       wg,
		taskChan: taskChan,
		doneChan: doneChan,
	}
}

func (w *worker) start() {
	go func(ctx context.Context, wg signaler, taskChan chan Runnable, doneChan chan bool) {
		for {
			select {
			case <-ctx.Done():
				wg.Done()
				return

			case task := <-taskChan:
				task.Run()
				doneChan <- true
			}
		}
	}(w.ctx, w.wg, w.taskChan, w.doneChan)
}
