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
	Tasks      []Task   `json:"tasks"`
	TaskID     string   `json:"task_id"`
	Error      string   `json:"error"`
	ExecTime   string   `json:"exec_time"`
	FinishTime string   `json:"finish_time"`
}

type Result struct {
	SubResultCode int    `json:"code"`
	TaskType      int    `json:"task_type"`
	SubResult     string `json:"sub_result"`
	SubTask       string `json:"sub_task"`
	SubTaskID     string `json:"sub_task_id"`
	Error         string `json:"error"`
}

func NewResult(task *Task) *Result {
	return &Result{
		TaskType:      task.TaskType,
		SubResultCode: CHECK_ERROR,
		SubTaskID:     task.SubTaskID,
	}
}

func NewGatewayResult(task *ServerTask) *ServerResult {
	result := &ServerResult{
		StatusCode: http.StatusOK,
		ExecTime:   time.Now().Format("2006-01-02 15:04:05"),
		Tasks:      task.Tasks,
	}
	result.TaskID = task.TaskID
	return result
}
