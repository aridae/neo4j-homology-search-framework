package commands

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	yaml "gopkg.in/yaml.v2"

	dataaccess "github.com/aridae/neo4j-homology-search-framework/client/v1/internal/data-access"
	db "github.com/aridae/neo4j-homology-search-framework/client/v1/internal/dbdriver"
	mdl "github.com/aridae/neo4j-homology-search-framework/client/v1/internal/model"
	"github.com/aridae/neo4j-homology-search-framework/client/v1/internal/utils/workerspool"
)

var (
	ParsingArgs = []string{
		"./preprocessing/v1/parse-fasta.sh",
		"-dout", "/fasta/jsons",
		"-mc", "1000000",
	}
)

type AddGenomeCommand struct {
	Header CommandHeader
	Data   []byte
}

var (
	pool *workerspool.WorkersPool = workerspool.GetWorkersPool(16, 100)
)

func AAAddSequence(repo dataaccess.Repository, seq *mdl.Sequence) error {
	log.Println("processing sequence", seq.Name)
	ctx := context.Background()
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		pool.RunBackground(ctx)
	}()

	// trim, split
	kmersStr := strings.Trim(seq.Data, "\n ")
	kmers := strings.Split(kmersStr, ",")
	for i, kmer := range kmers {
		val_cnt := strings.Split(kmer, ":")

		val := val_cnt[0]
		cnt, _ := strconv.ParseInt(val_cnt[1], 10, 64)
		pool.AddTask(NewAddGenomeTask(repo, seq, val, cnt))

		if i%1000 == 0 {
			log.Printf("Pushed %d kmers of sequence %s already\n", i, seq.Name)
		}
	}
	pool.Finish()
	wg.Wait()
	return nil
}

func AAAddGenome(Data []byte) error {
	// добавляем пул в ожидаемые горутины
	// важно: весь пул крутится в одной отдельной горутине на фоне
	// (ее мы и добавляем) - он сам следит за выделением и завершением
	// своих воркеров

	log.Printf("Adding genome - \n%s..., %d bytes\n", string(Data[:200]), len(Data))
	var genome mdl.Genome
	err := yaml.Unmarshal(Data, &genome)
	if err != nil {
		log.Println("failed to unmarshall genome data to genome struct:", err)
		return err
	}

	// add genome meta to database
	// создаем клиента для бд
	neo4jClient, err := db.GetNeo4jClient(&db.Options{
		URI:      "bolt://159.89.9.159:5001/",
		DB:       "neo4j",
		User:     "neo4j",
		Password: "H7rxhdt6!-jwt",
	})
	if err != nil {
		log.Println(err)
		return err
	}

	// перед выходом почистим
	defer neo4jClient.CloseNeo4jClient()
	defer log.Println("db client closed")

	repo := dataaccess.NewNeo4jRepository(neo4jClient)
	err = repo.AddGenomeMeta(&genome)
	if err != nil {
		log.Println("failed to add genome meta:", err)
		return err
	}

	// add sequence to database
	for i, seq := range genome.Sequences {
		if i == 4 {
			err = AAAddSequence(repo, &seq)
			if err != nil {
				log.Printf("failed to add sequence [%d] - %s: %s\n", i, seq.Name, err)
				return err
			}
		}
	}

	return nil
}

func NewAddGenomeCommand(path string, k int64) (Command, error) {
	ParsingArgs := append(ParsingArgs, "-fin", path, "-k", strconv.FormatInt(k, 10))

	fmt.Printf("%+v\n", ParsingArgs)

	// запускаем скрипт, дложидаемся завершения
	bash := exec.Command("/bin/bash", ParsingArgs...)

	stdout, _ := bash.StdoutPipe()
	stderr, _ := bash.StderrPipe()

	if err := bash.Start(); err != nil {
		log.Printf("Error executing command: %s......\n", err.Error())
		return nil, err
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		io.Copy(os.Stderr, stderr)
	}()

	wg.Add(1)
	// contentChan := make(chan []byte) // пусть чтение будет синхронным
	// errorChan := make(chan error)

	content := make([]byte, 0)
	go func(content *[]byte) {
		defer wg.Done()
		data, _ := io.ReadAll(stdout)
		(*content) = append((*content), data...)

		// // читаем из пайпа команды данные чанками по 4кб
		// // чанки прокидываем в канал, который читает основная горутина
		// // если эта горутина падает - надо указать основной горутине
		// // чистенько завершиться - перемешивать статусы с данными не хочется
		// // можно завести канал для статусов - либо структуру, где в заголовочнике
		// // будет статус, а в теле данные, но это вроде усложнение
		// r := bufio.NewReader(stdout)
		// buf := make([]byte, 0, 4<<20)
		// for {
		// 	n, err := r.Read(buf[:cap(buf)])
		// 	buf = buf[:n]
		// 	if err != nil {
		// 		if err == io.EOF {
		// 			break // дошли до конца пайпа
		// 		} else {
		// 			chanerr <- err // на проблемку напали
		// 			return
		// 		}
		// 	}
		// 	contentChan <- buf
		// }

		// data, _ := ioutil.ReadAll(stdout)
		// contentChan <- data
		// close(contentChan)
		// close(errorChan)
	}(&content)

	// // основная горутина читает из каналов и следит, чтобы все ок было
	// content := make([]byte, 0)
	// done := false

	// в чем мем: горутина с чтецом закрывает канал, когда закончила
	// а как нам об этом узнать, чтобы завершить цикл опроса?
	// c помощью оков - если не ок, значит с каналом что-то не так
	// значит, занулляем его (а зачем занулять? потому что закрытый канал
	// ВСЕГДА возвращает занчение при попытке чтения - дефолтное для типа канала)
	// для нас это значит, что если канал с данными закрылся, то в селекте его ветка
	// ВСЕГДА отрабатывает с nil и в другие ветки мы не попадем - закрытость канала
	// можно проверить оком
	// for !done {
	// 	select {
	// 	case err, ok := <-errorChan:
	// 		if !ok {
	// 			// если канал закрылся, занулим его
	// 			errorChan = nil
	// 		} else if err != nil {
	// 			// если ты упал, у тебя провал, ты решил, что шансов больше нет
	// 			return nil, err
	// 		}
	// 		// дождаться завершения горутин и выйти
	// 	case chunk, ok := <-contentChan:
	// 		// добавить чанк к нашим чанкам
	// 		if !ok {
	// 			// если канал закрылся, занулим его и завершим цикл опроса
	// 			contentChan = nil
	// 			done = true
	// 		} else {
	// 			content = append(content, chunk...)
	// 		}
	// 	}
	// }

	// log.Printf("%+v\n", *NewCommandHeader(AddGenome, int64(len(content))))

	wg.Wait()
	log.Printf("waited for all gouroutines...\n")
	AAAddGenome(content)

	// // возвращаем в обработчик, чтобы передали бэку
	// return &AddGenomeCommand{
	// 	Header: *NewCommandHeader(AddGenome, int64(len(content))),
	// 	Data:   content,
	// }, nil
	return nil, nil
}

func (cmd *AddGenomeCommand) GetCmd() CommandOption {
	return cmd.Header.Cmd
}

func (cmd *AddGenomeCommand) GetHeader() CommandHeader {
	return cmd.Header
}
func (cmd *AddGenomeCommand) MarshalBody() ([]byte, int64) {
	return cmd.Data, int64(len(cmd.Data))
}

func (cmd *AddGenomeCommand) UnmarshalBody(b []byte) error {
	return nil
}
