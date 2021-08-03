package register

import (
	"context"

	"github.com/go-magic/mid-server/task"
)

/*
Factory 工厂函数
*/
type Factory func() task.Tasker

/*
Register 注册中心
*/
type Register interface {
	Register(int, Factory)
	UnRegister(int)
	Tasker(ctx context.Context, taskType int) task.Tasker
	Exit()
}
