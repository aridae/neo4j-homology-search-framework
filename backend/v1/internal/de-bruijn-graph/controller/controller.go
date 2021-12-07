package controller

import (
	"fmt"
	"log"
	"sync"

	dac "github.com/aridae/neo4j-homology-search-framework/backend/v1/internal/de-bruijn-graph/data-access"
	gen "github.com/aridae/neo4j-homology-search-framework/backend/v1/internal/utils/kmers-generator"
	"github.com/aridae/neo4j-homology-search-framework/backend/v1/internal/utils/workerspool"
)

const (
	InitEmptyDBG ControllerTaskOption = iota
)

// контроллер выделяет воркеров под таски,
// следит за обращениями к внешним данным и тд
type Controller struct {
	repo dac.Repository
	pool *workerspool.WorkersPool
}

func NewController(repo dac.Repository, pool *workerspool.WorkersPool) *Controller {
	return &Controller{
		repo: repo,
		pool: pool,
	}
}

func (controller *Controller) RunTask(task ControllerTaskOption) error {
	switch task {
	case InitEmptyDBG:
		return controller.InitEmptyDBG()
	default:
		return fmt.Errorf("unsupported controller task")
	}
}

func (controller *Controller) InitEmptyDBG() error {
	log.Println("Constructing Empty De bruijn Graph")
	log.Printf("Generating kmers [k=%d]...\n", 11)
	WordsGenerator := gen.NewKMersGenerator(11, 300000)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		WordsGenerator.Generate()
	}()

	log.Println("Fetching generated words...")
	iterCnt := 0
	for word := range WordsGenerator.OutKMers {
		err := controller.repo.MergekPlus1Mer([]byte(word))
		if err != nil {
			log.Println(err)
		}
		iterCnt++
		if iterCnt%10000 == 0 {
			log.Printf("Generated %d kmers already...\n", iterCnt)
		}
	}
	wg.Wait()
	return nil
}
