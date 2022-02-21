package workerpool

import (
	"sync"
)

type WorkerPool struct {
	works chan func()
	wg    sync.WaitGroup
}

func New() *WorkerPool {
	wp := &WorkerPool{
		wg: sync.WaitGroup{},
	}
	return wp
}

func (w *WorkerPool) InitChan() {
	w.works = make(chan func())
}

func (w *WorkerPool) AddWorkers(workers int) {
	w.wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			for work := range w.works {
				work()
			}
			w.wg.Done()
		}()
	}
}

func (w *WorkerPool) Add(fn func()) {
	w.works <- fn
}

func (w *WorkerPool) ShutDown() {
	close(w.works)
	w.wg.Wait()
}
