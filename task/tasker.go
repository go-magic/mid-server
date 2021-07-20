package task

type TaskerFactory func() Tasker

type Tasker interface {
	Check(subTask string) (subResult string, err error)
}
