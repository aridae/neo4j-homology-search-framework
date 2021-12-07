package workerspool

// абстрактная задача
type Task interface {
	GetID() int     // идентифицировать задачу
	Process() error // выполнить задачу
	Cleanup() error // что сделать после того, как функция выполнится
}

func process(workerID int, task Task) {
	// log.Printf("Worker %d processes task %v\n", workerID, task.GetID())
	task.Process()
	task.Cleanup()
	// log.Printf("Worker %d did task %v\n", workerID, task.GetID())
}
