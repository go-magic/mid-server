package register

import (
	"context"

	"github.com/go-magic/mid-server/task"
)

type taskTypeFactory struct {
	taskType int
	factory  Factory
}

type registerCenter struct {
	center           map[int]Factory
	registerChan     chan taskTypeFactory
	cancellationChan chan int
	getTaskerChan    chan int
	sendTaskerChan   chan task.Tasker
	exist            chan struct{}
}

func NewRegisterCenter() *registerCenter {
	r := &registerCenter{}
	r.center = make(map[int]Factory)
	r.registerChan = make(chan taskTypeFactory)
	r.cancellationChan = make(chan int)
	r.getTaskerChan = make(chan int)
	r.sendTaskerChan = make(chan task.Tasker)
	r.exist = make(chan struct{})
	go r.start()
	return r
}

func (r registerCenter) start() {
	for {
		select {
		case msg := <-r.registerChan:
			r.register(msg)
		case cancel := <-r.cancellationChan:
			r.unRegister(cancel)
		case taskType := <-r.getTaskerChan:
			r.sendTasker(taskType)
		case <-r.exist:
			return
		}
	}
}

func (r registerCenter) Register(taskType int, factory Factory) {
	r.registerChan <- taskTypeFactory{taskType: taskType, factory: factory}
}

func (r registerCenter) register(msg taskTypeFactory) {
	r.center[msg.taskType] = msg.factory
}

func (r registerCenter) UnRegister(taskType int) {
	r.cancellationChan <- taskType
}

func (r registerCenter) unRegister(taskType int) {
	delete(r.center, taskType)
}

func (r registerCenter) Tasker(ctx context.Context, taskType int) task.Tasker {
	r.getTaskerChan <- taskType
	for {
		select {
		case <-ctx.Done():
			return nil
		case tasker := <-r.sendTaskerChan:
			return tasker
		}
	}
}

func (r registerCenter) sendTasker(taskType int) {
	factory := r.center[taskType]
	if factory == nil {
		r.sendTaskerChan <- nil
		return
	}
	r.sendTaskerChan <- factory()
}

func (r registerCenter) Exit() {
	r.exist <- struct{}{}
}
