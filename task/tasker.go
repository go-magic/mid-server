package task

/*
Tasker 执行任务接口,参数和返回值为json格式
*/
type Tasker interface {
	Check(subTask string) (subResult string, err error)
}
