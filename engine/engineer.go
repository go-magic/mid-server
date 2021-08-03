package engine

import "github.com/go-magic/mid-server/task"

/*
Producer 生产者,可以从网关获取任务,也可以来自路由任务,也可以自己创建任务。包括请求封装和拆分等，可以做装饰处理。
*/
type Producer interface {
	CreateServerTask() (*task.ServerTask, error)
}

/*
Consumer 消费者,可以由网关消费,也可以发送给前端,也可以自己内部消化。包括结果封装和拆分等，可以做装饰处理。
*/
type Consumer interface {
	DestroyServerResult(result *task.ServerResult) error
}
