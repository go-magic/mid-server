package dispatcher

import (
	"github.com/go-magic/mid-server/task"
	"github.com/go-magic/mid-server/work"
)

type Dispatch struct {
	maxWorkers        int
	workerPool        chan chan task.CheckRequest
	checkRequestQueue chan task.CheckRequest
	exitChan          chan struct{}
	workers           []work.Worker
}

func NewDispatcher(maxWorkers int) *Dispatch {
	pool := make(chan chan task.CheckRequest, maxWorkers)
	exit := make(chan struct{})
	checkRequestQueue := make(chan task.CheckRequest)
	d := &Dispatch{
		workerPool:        pool,
		maxWorkers:        maxWorkers,
		exitChan:          exit,
		checkRequestQueue: checkRequestQueue,
	}
	d.start()
	return d
}

func (d *Dispatch) start() {
	for i := 0; i < d.maxWorkers; i++ {
		worker := work.NewWorker(d.workerPool)
		worker.Start()
		d.workers = append(d.workers, worker)
	}
	go d.dispatch()
}

func (d *Dispatch) dispatch() {
	for {
		select {
		case request := <-d.checkRequestQueue:
			go d.addCheckRequest(request)
		case <-d.exitChan:
			return
		}
	}
}

func (d *Dispatch) AddExecuteTasker(tasker task.Tasker, subTask *task.Task, ch chan task.CheckResult) {
	d.checkRequestQueue <- task.CheckRequest{SubTask: subTask, Tasker: tasker, CheckResultChan: ch}
}

func (d *Dispatch) addCheckRequest(request task.CheckRequest) {
	subTaskChannel := <-d.workerPool
	subTaskChannel <- request
}

func (d *Dispatch) Exit() {
	d.exitChan <- struct{}{}
}

func (d *Dispatch) exit() {
	for _, worker := range d.workers {
		worker.Exit()
	}
}
