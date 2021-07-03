package schedule

import "mid-server/work"

type Pool struct {
	Workers chan work.Worker
}

func NewPool(maxRoutine int) *Pool {
	return &Pool{
		Workers: make(chan work.Worker, maxRoutine),
	}
}
