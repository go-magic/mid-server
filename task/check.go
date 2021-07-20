package task

import "net/http"

type CheckRequest struct {
	SubTask         *Task
	CheckResultChan chan CheckResult
	Tasker          Tasker
}

type CheckResult struct {
	SubResult *Result
	Error     error
}

func CreateCheckRequest(subTask *Task, checkResultChan chan CheckResult, tasker Tasker) CheckRequest {
	return CheckRequest{
		SubTask:         subTask,
		CheckResultChan: checkResultChan,
		Tasker:          tasker,
	}
}

func ErrorResult(subTask *Task, err string) *Result {
	result := NewResult(subTask)
	result.Error = err
	return result
}

func SuccessResult(subTask *Task, subResult string) *Result {
	result := NewResult(subTask)
	result.SubResultCode = http.StatusOK
	result.SubResult = subResult
	return result
}
