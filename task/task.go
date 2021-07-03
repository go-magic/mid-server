package task

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
