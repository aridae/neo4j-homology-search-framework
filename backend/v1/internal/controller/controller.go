package controller

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	cmd "github.com/aridae/neo4j-homology-search-framework/backend/v1/internal/commands"
	dac "github.com/aridae/neo4j-homology-search-framework/backend/v1/internal/data-access"
	"github.com/aridae/neo4j-homology-search-framework/backend/v1/internal/model"
	"github.com/aridae/neo4j-homology-search-framework/backend/v1/internal/utils/chunkreader"
	gen "github.com/aridae/neo4j-homology-search-framework/backend/v1/internal/utils/kmers-generator"
	"github.com/aridae/neo4j-homology-search-framework/backend/v1/internal/utils/workerspool"
)

type ControllerTaskOption int64

const (
	InitEmptyDBG ControllerTaskOption = iota

	ChunkSize          = 1024 * 4
	GeneratorQueueSize = 300000

	WorkersCnt       = 8
	WorkersQueueSize = 10
)

// контроллер принимает команды - слушает сокет
// некоторые команды дробятся на таски и выполняются параллельно
// некоторые последовательно целиком
// не все таски могут быть выполнены, пока работают другие таски, поэтому лочим
// TODO: мб ввести джобу или пайплайн - набор команд с указанием, как их выполнять
type Controller struct {
	repo dac.Repository
	pool *workerspool.WorkersPool

	// здесь будет статус модели
	// статус модели обновляется после выполнение команды
	// это задел для джобов и пайплайнов(?)
	Model *model.Model

	// тут абстрактная команда подменяется
	// указателем на реальную команду
	CommandsQueue chan cmd.Command
}

func NewController(repo dac.Repository, In chan cmd.Command) *Controller {
	// создаем пул воркеров, но не запускаем
	pool := workerspool.GetWorkersPool(WorkersCnt, WorkersQueueSize)
	return &Controller{
		repo:          repo,
		pool:          pool,
		CommandsQueue: In,
		Model:         model.EmptyModel(),
	}
}

// здесь мы инициализируем пул фоновых работяг,
// которые на фоне крутятся и ждут таски,
func (controller *Controller) RunControl(ctx context.Context) {

	log.Printf("Controller started\n")
	defer log.Println("Controller stopped")

	// контроллер завершается только после того, как завершится
	// пул воркеров, а пул воркеров завершится когда завершится последний работяга
	wg := &sync.WaitGroup{}

	// добавляем пул в ожидаемые горутины
	// важно: весь пул крутится в одной отдельной горутине на фоне
	// (ее мы и добавляем) - он сам следит за выделением и завершением
	// своих воркеров
	wg.Add(1)
	go func() {
		defer wg.Done()
		controller.pool.RunBackground(ctx)
	}()

	for {
		log.Printf("controller loop\n")
		select {
		case <-ctx.Done():
			// завершаем работяг и сами завершаемся
			log.Println("finishing controller...")

			controller.pool.Finish()
			wg.Wait()
			return
		case command := <-controller.CommandsQueue:
			// берем команду и отдаем пулу на вылонение
			log.Println("fetched new command!")
			controller.handleCommand(command)
		}
	}
}

func (controller *Controller) handleCommand(command cmd.Command) error {
	switch command.GetHeader().Cmd {
	case cmd.InitEmptyGraph:
		cmdInitEmptyGraph, _ := command.(*cmd.InitEmptyGraphCommand)
		return controller.InitEmptyDBG(cmdInitEmptyGraph.K)
	case cmd.AddGenome:
		cmdAddGenome, _ := command.(*cmd.AddGenomeCommand)
		return controller.AddGenome(cmdAddGenome.Path)
	default:
		return fmt.Errorf("unsupported controller command")
	}
}

func (controller *Controller) InitEmptyDBG(k int64) error {
	log.Println("Constructing Empty De bruijn Graph")
	log.Printf("Generating kmers [k=%d]...\n", k)

	// обновляем данные о модели
	controller.Model.SetK(k)

	// создаем генератор
	WordsGenerator := gen.NewKMersGenerator(uint64(k+1), GeneratorQueueSize)

	// в отдельном потоке генератор генерит кмеры и
	// пишет их в канал
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		WordsGenerator.Generate()
	}()

	// основной поток читает кмеры из канала и отдает воркерам замержить в бд
	// надо подключить кафку, чтобы мержили воркеры, а не один поток
	log.Println("Fetching generated k+1-mers...")
	iterCnt := 0
	for word := range WordsGenerator.OutKMers {
		// TODO: добавить кафку - только тут траблы с тем как партицировать графовую бд
		// controller.pool.AddTask(
		// 	NewInitGraphTask(
		// 		[]byte(word),
		// 		controller.repo,
		// 	),
		// )
		controller.repo.MergePrecedingKMers([]byte(word))
		iterCnt++
		if iterCnt%10000 == 0 {
			log.Printf("Generated %d k+1-mers already...\n", iterCnt)
		}
	}
	log.Printf("Generated %d k+1-mers. The base is ready!\n", iterCnt)
	wg.Wait()
	return nil
}

// еще возникла проблема, в коммьюнити едишн нео4ж можно только одну базу на сервер
// а чтобы тестировать, нужно будет две базы - игрушечная и нормальная
// мб можно завести два сервиса с докеровскими нео4ж и смапить их на разные порты
// мб можно и приложеньку сбилдить как докеровский сервис и деплоить ее системд сервисом
// проблема только в том, что собираться будет дольше
func (controller *Controller) AddGenome(path string) error {
	log.Printf("Adding genome, path=%s\n", path)

	// открываем файл
	fasta, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return err
	}
	defer fasta.Close()

	// создаем ридера
	reader := chunkreader.GetChunkReader(ChunkSize, int(controller.Model.K), fasta)

	// в основном потоке ридер читает чанки памяти и
	// отдает их воркерам(!)
	chunk, err := reader.ReadChunk()
	for err == nil {
		//controller.repo.MergePrecedingKMers()
		// TODO: добавить кафку - только тут траблы с тем как партицировать графовую бд
		controller.pool.AddTask(
			NewAddGenomeTask(
				chunk,
				controller.repo,
				reader,
			),
		)
		chunk, err = reader.ReadChunk()
	}
	if err != io.EOF {
		log.Println(err)
		return err
	}

	return nil
}
