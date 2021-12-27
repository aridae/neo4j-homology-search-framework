package workerspool

type Worker struct {
	ID         int
	tasksQueue chan Task // работяга работает, пока есть работа :_)
}

func NewWorker(tasksQueue chan Task, ID int) *Worker {
	return &Worker{
		ID:         ID,
		tasksQueue: tasksQueue,
	}
}

func (w *Worker) RunBackground() {
	for task := range w.tasksQueue {
		task.Process()
		task.Cleanup()
	}
}
