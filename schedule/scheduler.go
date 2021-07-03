package schedule

import (
	"context"
	"mid-server/task"
)

type Scheduler interface {
	Execute(context.Context, []task.Task) ([]task.Result, error)
}
