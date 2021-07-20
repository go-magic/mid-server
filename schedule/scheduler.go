package schedule

import (
	"context"

	"github.com/go-magic/mid-server/task"
)

type Scheduler interface {
	Execute(context.Context, []task.Task) ([]task.Result, error)
	Exit()
}
