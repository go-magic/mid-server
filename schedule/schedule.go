package schedule

import (
	"context"
	"errors"
	"mid-server/register"
	"mid-server/task"
	"mid-server/work"
	"net/http"
	"sync"
)

const (
	MIN_VALID_ROUTINE = 1
	MAX_VALID_ROUTINE = 100000
)

type checkResult struct {
	Result *task.Result
	E      error
}

type schedule struct {
	resultChan chan checkResult
	register   register.Register
	workers    chan work.Worker
}

func NewSchedule(maxRoutine int, r register.Register) *schedule {
	if maxRoutine < MIN_VALID_ROUTINE || maxRoutine > MAX_VALID_ROUTINE {
		return nil
	}
	s := &schedule{}
	s.resultChan = make(chan checkResult)
	s.register = r
	s.workers = make(chan work.Worker, maxRoutine)
	return s
}

func (s *schedule) Execute(ctx context.Context, subTasks []task.Task) ([]task.Result, error) {
	results := make([]task.Result, 0, len(subTasks))
	group := sync.WaitGroup{}
	checkoutChan := make(chan struct{})
	resultChan := make(chan checkResult)
	for _, subTask := range subTasks {
		group.Add(1)
		go s.check(subTask, resultChan)
	}
	go func() {
		group.Wait()
		checkoutChan <- struct{}{}
	}()
	for {
		select {
		case result := <-resultChan:
			results = append(results, *result.Result)
			group.Done()
		case <-checkoutChan:
			return results, nil
		case <-ctx.Done():
			return results, errors.New("time out")
		}
	}
}

func (s *schedule) check(subTask task.Task, resultChan chan checkResult) {
	t := s.register.Tasker(context.Background(), subTask.TaskType)
	if t == nil {
		resultChan <- checkResult{
			Result: task.NewResult(&subTask),
			E:      errors.New("messageID not register"),
		}
		return
	}
	result, err := s.checkByTemplate(&subTask, t)
	resultChan <- checkResult{
		Result: result,
		E:      err,
	}
}

func (s *schedule) checkByTemplate(task *task.Task, tasker task.Tasker) (*task.Result, error) {
	subResult, checkErr := tasker.Check(task)
	if checkErr != nil {
		return s.errorResult(task, checkErr.Error()), checkErr
	}
	return s.successResult(task, subResult), nil
}

func (s *schedule) errorResult(subTask *task.Task, err string) *task.Result {
	result := task.NewResult(subTask)
	result.Error = err
	return result
}

func (s *schedule) successResult(subTask *task.Task, subResult string) *task.Result {
	result := task.NewResult(subTask)
	result.SubResultCode = http.StatusOK
	result.SubResult = subResult
	return result
}
