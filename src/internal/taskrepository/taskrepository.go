package taskrepository

type Taskrepository interface {
	AddTask(string, int) error
}
