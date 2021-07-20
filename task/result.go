package task

import (
	"net/http"
	"time"
)

const (
	CHECK_ERROR         = 1000
	CHECK_TIMEOUT_ERROR = 1001
)

type ServerResult struct {
	StatusCode int      `json:"status_code"`
	Results    []Result `json:"results"`
	TaskID     string   `json:"task_id"`
	Error      string   `json:"error"`
	ExecTime   string   `json:"exec_time"`
	FinishTime string   `json:"finish_time"`
}

type Result struct {
	SubResultCode int    `json:"code"`
	TaskType      string `json:"task_type"`
	SubResult     string `json:"result"`
	SubTaskID     string `json:"sub_task_id"`
	Error         string `json:"error"`
	Task          Task   `json:"task"`
}

func NewResult(task *Task) *Result {
	return &Result{
		TaskType:      task.TaskType,
		SubResultCode: CHECK_ERROR,
		SubTaskID:     task.SubTaskID,
		Task:          *task,
	}
}

func NewGatewayResult(task *ServerTask) *ServerResult {
	result := &ServerResult{
		StatusCode: http.StatusOK,
		ExecTime:   time.Now().Format("2006-01-02 15:04:05"),
	}
	result.TaskID = task.TaskID
	return result
}
