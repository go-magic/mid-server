package register

import (
	"context"

	"github.com/go-magic/mid-server/task"
)

type Factory func() task.Tasker

type Register interface {
	Register(string, Factory)
	UnRegister(string)
	Tasker(ctx context.Context, taskType string) task.Tasker
	Exit()
}
