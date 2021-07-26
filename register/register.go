package register

import (
	"context"

	"github.com/go-magic/mid-server/task"
)

type Factory func() task.Tasker

type Register interface {
	Register(int, Factory)
	UnRegister(int)
	Tasker(ctx context.Context, taskType int) task.Tasker
	Exit()
}
