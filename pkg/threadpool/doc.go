// Package threadpool offers a pool of workers implemented using goroutines.
//
// Create a new thread pool by calling NewThreadPool and specify the number of workers wanted and a work queue size.
// If the work queue is full when a new task is pushed, the thread pool will block until another task finishes.
//
// Tasks can be any objects that implement the Runnable interface.
package threadpool
