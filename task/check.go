package task

import "net/http"

/*
CheckRequest 任务执行请求参数
*/
type CheckRequest struct {
	SubTask         *Task
	CheckResultChan chan CheckResult
	Tasker          Tasker
}

/*
CheckResult 任务执行完毕结果回传结构
*/
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

/*
ErrorResult 返回错误结果
*/
func ErrorResult(subTask *Task, err string) *Result {
	result := NewResult(subTask)
	result.Error = err
	return result
}

/*
SuccessResult 返回成功结果
*/
func SuccessResult(subTask *Task, subResult string) *Result {
	result := NewResult(subTask)
	result.SubResultCode = http.StatusOK
	result.SubResult = subResult
	return result
}
