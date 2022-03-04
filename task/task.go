package task

type ServerTask struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
	TaskID  string `json:"taskId"`
	Tasks   []Task `json:"tasks"`
}

// Task subTask
type Task struct {
	TaskType  int `json:"taskType"`
	SubTask   string `json:"subTask"`
	SubTaskID string `json:"subTaskId"`
}
