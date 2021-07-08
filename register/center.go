package register

import (
	"context"

	"github.com/go-magic/mid-server/task"
)

type taskTypeFactory struct {
	taskType string
	factory  Factory
}

type registerCenter struct {
	center           map[string]Factory
	registerChan     chan taskTypeFactory
	cancellationChan chan string
	getTaskerChan    chan string
	sendTaskerChan   chan task.Tasker
	exist            chan struct{}
}

func NewRegisterCenter() *registerCenter {
	r := &registerCenter{}
	r.center = make(map[string]Factory)
	r.registerChan = make(chan taskTypeFactory)
	r.cancellationChan = make(chan string)
	r.getTaskerChan = make(chan string)
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

func (r registerCenter) Register(taskType string, factory Factory) {
	r.registerChan <- taskTypeFactory{taskType: taskType, factory: factory}
}

func (r registerCenter) register(msg taskTypeFactory) {
	r.center[msg.taskType] = msg.factory
}

func (r registerCenter) UnRegister(taskType string) {
	r.cancellationChan <- taskType
}

func (r registerCenter) unRegister(taskType string) {
	delete(r.center, taskType)
}

func (r registerCenter) Tasker(ctx context.Context, taskType string) task.Tasker {
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

func (r registerCenter) sendTasker(taskType string) {
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
