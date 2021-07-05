package work

import "github.com/go-magic/mid-server/task"

type Worker interface {
	Do(task *task.Task) (string, error)
}
