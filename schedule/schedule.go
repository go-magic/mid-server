package schedule

import (
	"context"
	"errors"
	"sync"

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

func NewSchedule(r register.Register, dis dispatcher.Dispatcher) *schedule {
	if r == nil {
		return nil
	}
	if dis == nil {
		dis = &dispatcher.Dispatch{}
	}
	s := &schedule{}
	s.dis = dis
	s.register = r
	return s
}

func (s schedule) Execute(ctx context.Context, subTasks []task.Task) ([]task.Result, error) {
	results := make([]task.Result, 0, len(subTasks))
	group := sync.WaitGroup{}
	checkoutChan := make(chan struct{}, len(subTasks))
	resultChan := make(chan task.CheckResult)
	group.Add(len(subTasks))
	go s.checkResults(subTasks, resultChan)
	go func() {
		group.Wait()
		checkoutChan <- struct{}{}
	}()
	for {
		select {
		case result := <-resultChan:
			results = append(results, *result.SubResult)
			group.Done()
		case <-checkoutChan:
			return results, nil
		case <-ctx.Done():
			group.Done()
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
