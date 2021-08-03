package example

import (
	"github.com/go-magic/mid-server/dispatcher"
	"github.com/go-magic/mid-server/engine"
	"github.com/go-magic/mid-server/register"
	"github.com/go-magic/mid-server/schedule"
	"github.com/go-magic/mid-server/task"
	"log"
	"net/http"
	"time"
)

const (
	MSG_HTTP_TASK = iota + 1
)

const (
	MAX_ROUTIME = 10
)

type Create struct {
}

func NewCreate() *Create {
	return &Create{}
}

func NewServerTask() *task.ServerTask {
	serverTask := task.ServerTask{}
	serverTask.Tasks = initTasks()
	serverTask.TaskID = "1"
	serverTask.Code = http.StatusOK
	serverTask.Status = "success"
	serverTask.Message = "请求成功"
	return &serverTask
}

func (c Create) CreateServerTask() (*task.ServerTask, error) {
	return NewServerTask(), nil
}

type Destroy struct {
}

func NewDestroy() *Destroy {
	return &Destroy{}
}

func (d Destroy) DestroyServerResult(result *task.ServerResult) error {
	log.Println("执行完毕", result)
	return nil
}

/*
初始化engine config
*/
func initEngineConfig() engine.Config {
	producer := NewCreate()
	consumer := NewDestroy()
	c := engine.Config{}
	c.ErrorWaitTime = time.Second * 1
	c.SuccessWaitTime = time.Second * 1
	c.ExecuteTime = time.Second * 1
	c.Producer = producer
	c.Consumer = consumer
	return c
}

/*
初始化engine
*/
func initEngine(scheduler schedule.Scheduler) error {
	engine.GetEngineInstance().SetScheduler(scheduler)
	engine.GetEngineInstance().SetConfig(initEngineConfig())
	return nil
}

func initScheduler(register register.Register, maxRoutine int) schedule.Scheduler {
	return schedule.NewSchedule(register, dispatcher.NewDispatcher(maxRoutine))
}

func initRegister() register.Register {
	r := register.NewRegisterCenter()
	r.Register(MSG_HTTP_TASK, NewHttpServer)
	return r
}

/*
初始化服务
*/
func initServer() error {
	r := initRegister()
	scheduler := initScheduler(r, MAX_ROUTIME)
	if err := initEngine(scheduler); err != nil {
		return err
	}
	return engine.GetEngineInstance().Start()
}

/*
开始服务
*/
func startServer() error {
	return initServer()
}
