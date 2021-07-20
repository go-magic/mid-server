package schedule

import (
	"context"
	"errors"
	"github.com/go-magic/mid-server/dispatcher"
	"sync"

	"github.com/go-magic/mid-server/register"
	"github.com/go-magic/mid-server/task"
)

type schedule struct {
	register register.Register
	//限流
	dis *dispatcher.Dispatcher
}

func NewSchedule(maxGoRoutine int, r register.Register) *schedule {
	if r == nil {
		return nil
	}
	dis := dispatcher.NewDispatcher(maxGoRoutine)
	s := &schedule{}
	s.dis = dis
	s.register = r
	dis.Run()
	return s
}

func (s schedule) Execute(ctx context.Context, subTasks []task.Task) ([]task.Result, error) {
	results := make([]task.Result, 0, len(subTasks))
	group := sync.WaitGroup{}
	checkoutChan := make(chan struct{})
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
			return results, errors.New("time out")
		}
	}
}

func (s schedule) checkResults(subTasks []task.Task, resultChan chan task.CheckResult) {
	for i, subTask := range subTasks {
		tasker := s.register.Tasker(context.Background(), subTask.TaskType)
		requestChan := task.CreateCheckRequest(&subTasks[i], resultChan, tasker)
		s.dis.AddCheckRequest(requestChan)
	}
}

func (s schedule) Check(task *task.Task, tasker task.Tasker) (*task.Result, error) {
	return tasker.Check(task)
}
