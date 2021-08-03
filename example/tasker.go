package example

import (
	"encoding/json"
	"github.com/go-magic/mid-server/task"
	"net/http"
)

type HttpTask struct {
	Url string `json:"url"`
}

type HttpResult struct {
	StatusCode int `json:"url"`
}

type HttpServer struct {
	subTask HttpTask
}

func NewHttpServer() task.Tasker {
	return &HttpServer{}
}

func (h *HttpServer) parse(subTask string) error {
	return json.Unmarshal([]byte(subTask), &h.subTask)
}

func (h *HttpServer) Check(subTask string) (string, error) {
	err := h.parse(subTask)
	if err != nil {
		return "", err
	}
	client := http.Client{}
	req, _ := http.NewRequest("GET", h.subTask.Url, nil)
	res, doErr := client.Do(req)
	if doErr != nil {
		return "", doErr
	}
	result := HttpResult{}
	result.StatusCode = res.StatusCode
	b, marshalErr := json.Marshal(result)
	if marshalErr != nil {
		return "", marshalErr
	}
	return string(b), nil
}

func initTasks() []task.Task {
	tasks := make([]task.Task, 0)
	subTask := task.Task{}
	httpTask := HttpTask{Url: "https://www.baidu.com"}
	b, err := json.Marshal(httpTask)
	if err != nil {
		return tasks
	}
	subTask.SubTaskID = "1"
	subTask.SubTask = string(b)
	subTask.TaskType = MSG_HTTP_TASK
	tasks = append(tasks, subTask)
	return tasks
}
