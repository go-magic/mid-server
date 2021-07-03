package register

import (
	"context"
	"mid-server/task"
)

type Factory func() task.Tasker

type Register interface {
	Register(string, Factory)
	UnRegister(string)
	Tasker(ctx context.Context, taskType string) task.Tasker
	Exit()
}
