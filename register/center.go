package register

import (
	"context"

	"github.com/go-magic/mid-server/task"
)

type FactoryRegister struct {
	taskType int
	factory  Factory
}

type registerCenter struct {
	center           map[int]Factory      /*注册中心*/
	registerChan     chan FactoryRegister /*注册对象*/
	cancellationChan chan int             /*取消注册*/
	getTaskerChan    chan int             /*获取tasker*/
	sendTaskerChan   chan task.Tasker     /*发送tasker*/
	exit             chan struct{}        /*退出*/
}

func NewRegisterCenter() *registerCenter {
	r := &registerCenter{}
	r.center = make(map[int]Factory)
	r.registerChan = make(chan FactoryRegister)
	r.cancellationChan = make(chan int)
	r.getTaskerChan = make(chan int)
	r.sendTaskerChan = make(chan task.Tasker)
	r.exit = make(chan struct{})
	go r.start()
	return r
}

/*
启动服务
*/
func (r registerCenter) start() {
	for {
		select {
		case msg := <-r.registerChan:
			r.register(msg)
		case cancel := <-r.cancellationChan:
			r.unRegister(cancel)
		case taskType := <-r.getTaskerChan:
			r.sendTasker(taskType)
		case <-r.exit:
			return
		}
	}
}

/*
Register 外部注册接口
*/
func (r registerCenter) Register(taskType int, factory Factory) {
	r.registerChan <- FactoryRegister{taskType: taskType, factory: factory}
}

/*
内部注册接口
*/
func (r registerCenter) register(msg FactoryRegister) {
	r.center[msg.taskType] = msg.factory
}

/*
UnRegister 外部注销接口
*/
func (r registerCenter) UnRegister(taskType int) {
	r.cancellationChan <- taskType
}

/*
内部注销接口
*/
func (r registerCenter) unRegister(taskType int) {
	delete(r.center, taskType)
}

/*
Tasker 获取Tasker接口
*/
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

/*
内部发送tasker接口
*/
func (r registerCenter) sendTasker(taskType int) {
	factory := r.center[taskType]
	if factory == nil {
		r.sendTaskerChan <- nil
		return
	}
	r.sendTaskerChan <- factory()
}

/*
Exit 退出接口
*/
func (r registerCenter) Exit() {
	r.exit <- struct{}{}
}
