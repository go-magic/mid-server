package dispatcher

import (
	"github.com/go-magic/mid-server/task"
	"github.com/go-magic/mid-server/work"
)

type Dispatcher struct {
	maxWorkers        int
	workerPool        chan chan task.CheckRequest
	checkRequestQueue chan task.CheckRequest
	exitChan          chan struct{}
	workers           []work.Worker
}

func NewDispatcher(maxWorkers int) *Dispatcher {
	pool := make(chan chan task.CheckRequest, maxWorkers)
	exit := make(chan struct{})
	checkRequestQueue := make(chan task.CheckRequest)
	return &Dispatcher{
		workerPool:        pool,
		maxWorkers:        maxWorkers,
		exitChan:          exit,
		checkRequestQueue: checkRequestQueue,
	}
}

func (d *Dispatcher) Run() {
	for i := 0; i < d.maxWorkers; i++ {
		worker := work.NewWorker(d.workerPool)
		worker.Start()
		d.workers = append(d.workers, worker)
	}
	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case request := <-d.checkRequestQueue:
			go d.addCheckRequest(request)
		case <-d.exitChan:

		}
	}
}

func (d *Dispatcher) AddCheckRequest(request task.CheckRequest) {
	d.checkRequestQueue <- request
}

func (d *Dispatcher) addCheckRequest(request task.CheckRequest) {
	subTaskChannel := <-d.workerPool
	subTaskChannel <- request
}

func (d *Dispatcher) Exit() {
	d.exitChan <- struct{}{}
}

func (d *Dispatcher) exit() {
	for _, worker := range d.workers {
		worker.Stop()
	}
}
