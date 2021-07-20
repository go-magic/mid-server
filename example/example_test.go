package example

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/go-magic/mid-server/register"
	"github.com/go-magic/mid-server/schedule"
	"github.com/go-magic/mid-server/task"
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

func NewHttpGet() task.Tasker {
	return &HttpServer{}
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
	subTask.TaskType = "http_get"
	tasks = append(tasks, subTask)
	return tasks
}

func TestNewServer(t *testing.T) {
	center := register.NewRegisterCenter()
	center.Register("http_get", NewHttpGet)
	sch := schedule.NewSchedule(1, center)
	subResult, err := sch.Execute(context.Background(), initTasks())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(subResult)
}
