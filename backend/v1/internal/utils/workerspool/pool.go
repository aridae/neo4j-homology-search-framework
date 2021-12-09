package workerspool

import (
	"context"
	"log"
	"sync"
)

type WorkersPool struct {
	workersNumber int
	tasksQueue    chan Task
}

var (
	pool *WorkersPool
)

func GetWorkersPool(workersNumber, taskQueueSize int) *WorkersPool {
	if pool == nil {
		pool = &WorkersPool{
			workersNumber: workersNumber,
			tasksQueue:    make(chan Task, taskQueueSize),
		}
	}
	return pool
}

func (p *WorkersPool) AddTask(task Task) {
	p.tasksQueue <- task
}

func (p *WorkersPool) Finish() {
	// мы закроем канал с тасками так, что работяги дочитают все, что есть и завершатся
	close(p.tasksQueue)
}

func (p *WorkersPool) RunBackground(cxt context.Context) {
	log.Println("Initing rabotyag...")
	wg := &sync.WaitGroup{}
	wg.Add(p.workersNumber)

	// запускаем воркеров, которые будут
	// на фоне ждать таски
	// и вычитывать их из одного(!) канала с тасками
	// закрытый канал значит, что новых таск не будет
	// следовательно, заканчиваем все в очереди и завершаемся, работяги
	for i := 1; i <= p.workersNumber; i++ {
		worker := NewWorker(p.tasksQueue, i)
		go func() {
			defer wg.Done()
			worker.RunBackground()
		}()
	}

	// считаем, что пул работает, пока работают работяги
	wg.Wait()
	log.Println("All workers done")
}
