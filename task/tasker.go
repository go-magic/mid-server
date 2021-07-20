package task

type Tasker interface {
	Check(subTask string) (subResult string, err error)
}
