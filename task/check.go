package task

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
