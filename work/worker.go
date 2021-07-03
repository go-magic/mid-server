package work

import "mid-server/task"

type Worker interface {
	Do(task *task.Task) (string, error)
}
