package task

import (
	"net/http"
	"time"
)

const (
	CHECK_ERROR         = 1000
	CHECK_TIMEOUT_ERROR = 1001
)

/*
ServerResult 通用多任务结果
*/
type ServerResult struct {
	StatusCode int      `json:"statusCode"`
	Results    []Result `json:"results"`
	Tasks      []Task   `json:"tasks"`
	TaskID     string   `json:"taskId"`
	Error      string   `json:"error"`
	ExecTime   string   `json:"execTime"`
	FinishTime string   `json:"finishTime"`
}

/*
Result 通用单任务结果
*/
type Result struct {
	SubResultCode int    `json:"code"`
	TaskType      int    `json:"taskType"`
	SubResult     string `json:"subResult"`
	SubTask       string `json:"subTask"`
	SubTaskID     string `json:"subTaskId"`
	Error         string `json:"error"`
}

/*
NewResult 创建结果
*/
func NewResult(task *Task) *Result {
	return &Result{
		TaskType:      task.TaskType,
		SubResultCode: CHECK_ERROR,
		SubTaskID:     task.SubTaskID,
		SubTask:       task.SubTask,
	}
}

func NewServerResult(task *ServerTask) *ServerResult {
	result := &ServerResult{
		StatusCode: http.StatusOK,
		ExecTime:   time.Now().Format("2006-01-02 15:04:05"),
		Tasks:      task.Tasks,
	}
	result.TaskID = task.TaskID
	return result
}
