package task

type ServerTask struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
	TaskID  string `json:"task_id"`
	Tasks   []Task `json:"tasks"`
}

// Task subTask
type Task struct {
	TaskType  int `json:"task_type"`
	SubTask   string `json:"sub_task"`
	SubTaskID string `json:"sub_task_id"`
}
