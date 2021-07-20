package task

import "net/http"

type ServerTask struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Tasks   []Task `json:"tasks"`
	TaskID  string `json:"task_id"`
}

// Task subTask
type Task struct {
	TaskType  string `json:"task_type"`
	SubTask   string `json:"sub_task"`
	SubTaskID string `json:"sub_task_id"`
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
