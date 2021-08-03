package dispatcher

import "github.com/go-magic/mid-server/task"

type Dispatcher interface {
	AddExecuteTasker(tasker task.Tasker, subTask *task.Task, ch chan task.CheckResult)
	Exit()
}
