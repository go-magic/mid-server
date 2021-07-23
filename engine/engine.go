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

func (e *Engine) SetScheduler(scheduler schedule.Scheduler) {
	e.scheduler = scheduler
}

func (e *Engine) SetConfig(config Config) {
	e.config = config
}

func (e *Engine) start() {
	for {
		err := e.startServer()
		if err != nil {
			time.Sleep(e.config.ErrorWaitTime)
			continue
		}
		time.Sleep(e.config.SuccessWaitTime)
	}
}

func (e *Engine) init() error {
	if e.scheduler == nil ||
		e.config.Producer == nil ||
		e.config.Consumer == nil {
		return errors.New("args error")
	}
	if e.config.ErrorWaitTime == 0 {
		e.config.ErrorWaitTime = time.Second * 60
	}
	if e.config.SuccessWaitTime == 0 {
		e.config.SuccessWaitTime = time.Second * 60
	}
	if e.config.ExecuteTime == 0 {
		e.config.ExecuteTime = time.Second * 60
	}
	return nil
}

func (e *Engine) startServer() error {
	serverTask, err := e.config.Producer.CreateServerTask()
	if err != nil {
		return err
	}
	result, executeErr := e.execute(serverTask)
	if executeErr != nil {
		if result != nil {
			result.Error = executeErr.Error()
		}
	}
	return e.config.Consumer.DestroyServerResult(result)
}

func (e *Engine) Run() error {
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

func (e *Engine) execute(serverTask *task.ServerTask) (*task.ServerResult, error) {
	serverResult := task.NewGatewayResult(serverTask)
	ctx, _ := context.WithTimeout(context.Background(), e.config.ExecuteTime)
	results, err := e.scheduler.Execute(ctx, serverTask.Tasks)
	if err != nil {
		return serverResult, err
	}
	serverResult.Results = results
	return serverResult, nil
}

func (e *Engine) Exit() {
	e.scheduler.Exit()
	e.exit <- struct{}{}
}
