package task

type Tasker interface {
	Check(task *Task) (subResult string, err error)
}

