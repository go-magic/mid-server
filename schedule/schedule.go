package schedule

import (
	"context"
	"errors"
	"github.com/go-magic/mid-server/dispatcher"
	"github.com/go-magic/mid-server/register"
	"github.com/go-magic/mid-server/task"
)

/*
调度实体
*/
type schedule struct {
	register register.Register
	//限流
	dis dispatcher.Dispatcher
}

func NewSchedule(r register.Register, dis dispatcher.Dispatcher) schedule {
	s := schedule{}
	if r == nil {
		s.register = register.NewRegisterCenter()
		return s
	}
	if dis == nil {
		dis = &dispatcher.Dispatch{}
	}
	s.dis = dis
	s.register = r
	return s
}

func (s schedule) Execute(ctx context.Context, subTasks []task.Task) ([]task.Result, error) {
	results := make([]task.Result, 0, len(subTasks))
	if len(subTasks) == 0 {
		return results, nil
	}
	resultChan := make(chan task.CheckResult, len(subTasks))
	go s.checkResults(subTasks, resultChan)
	for {
		select {
		case result := <-resultChan:
			results = append(results, *result.SubResult)
			if len(results) == len(subTasks) {
				return results, nil
			}
		case <-ctx.Done():
			return results, errors.New("time out")
		}
	}
}

func (s schedule) checkResults(subTasks []task.Task, resultChan chan task.CheckResult) {
	for i, subTask := range subTasks {
		tasker := s.register.Tasker(context.Background(), subTask.TaskType)
		s.dis.AddExecuteTasker(tasker, &subTasks[i], resultChan)
	}
}

func (s schedule) Exit() {
	go s.register.Exit()
	go s.dis.Exit()
}
