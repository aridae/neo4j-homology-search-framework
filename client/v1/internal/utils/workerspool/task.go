package workerspool

// абстрактная задача
type Task interface {
	GetID() int     // идентифицировать задачу
	Process() error // выполнить задачу
	Cleanup() error // что сделать после того, как функция выполнится
}
