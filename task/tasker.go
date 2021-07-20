package task

type TaskerFactory func() Tasker

type Tasker interface {
	Check(task *Task) (result *Result, err error)
}
