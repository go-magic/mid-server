package schedule

import (
	"context"

	"github.com/go-magic/mid-server/task"
)

/*
Scheduler 调度接口
*/
type Scheduler interface {
	Execute(context.Context, []task.Task) ([]task.Result, error)
	Exit()
}
