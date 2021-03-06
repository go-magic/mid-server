package engine

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/go-magic/mid-server/schedule"
	"github.com/go-magic/mid-server/task"
)

type Engine struct {
	config    Config
	scheduler schedule.Scheduler
	exit      chan struct{}
	parser    ErrorParser
}

var (
	engine *Engine
	once   sync.Once
)

func GetEngineInstance() *Engine {
	once.Do(func() {
		engine = NewEngine()
	})
	return engine
}

func NewEngine() *Engine {
	return &Engine{
		exit: make(chan struct{}),
	}
}

/*
SetScheduler 设置调度
*/
func (e *Engine) SetScheduler(scheduler schedule.Scheduler) {
	e.scheduler = scheduler
}

func (e *Engine) SetParser(parser ErrorParser) {
	e.parser = parser
}

/*
SetConfig 设置engine配置文件
*/
func (e *Engine) SetConfig(config Config) {
	e.config = config
}

func (e *Engine) start() {
	for {
		err := e.startServer()
		if e.parser != nil {
			e.parser.Parser(err)
		}
	}
}

/*
初始化engine
*/
func (e *Engine) init() error {
	if e.scheduler == nil ||
		e.config.Producer == nil ||
		e.config.Consumer == nil {
		return errors.New("args error")
	}
	if e.config.ExecuteTime == 0 {
		e.config.ExecuteTime = time.Second * 60
	}
	return nil
}

/*
开启服务
*/
func (e *Engine) startServer() error {
	serverTask, err := e.config.Producer.CreateServerTask()
	if err != nil {
		return err
	}
	result, executeErr := e.execute(serverTask)
	if executeErr != nil {
		if result != nil {
			result.Error = executeErr.Error()
			result.StatusCode = task.CHECK_ERROR
		}
	}
	return e.config.Consumer.DestroyServerResult(result)
}

/*
Start 开始服务
*/
func (e *Engine) Start() error {
	if err := e.init(); err != nil {
		return err
	}
	go e.start()
	for {
		select {
		case <-e.exit:
			return nil
		}
	}
}

/*
调度执行多任务实例
*/
func (e *Engine) execute(serverTask *task.ServerTask) (*task.ServerResult, error) {
	serverResult := task.NewServerResult(serverTask)
	ctx, _ := context.WithTimeout(context.Background(), e.config.ExecuteTime)
	results, err := e.scheduler.Execute(ctx, serverTask.Tasks)
	serverResult.Results = results
	return serverResult, err
}

/*
Exit 退出
*/
func (e *Engine) Exit() {
	e.scheduler.Exit()
	e.exit <- struct{}{}
}
