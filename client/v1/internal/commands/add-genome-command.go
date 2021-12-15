package commands

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"sync"
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
	contentChan := make(chan []byte) // пусть чтение будет синхронным
	errorChan := make(chan error)
	go func(chanout chan []byte, chanerr chan error) {
		defer wg.Done()

		// читаем из пайпа команды данные чанками по 4кб
		// чанки прокидываем в канал, который читает основная горутина
		// если эта горутина падает - надо указать основной горутине
		// чистенько завершиться - перемешивать статусы с данными не хочется
		// можно завести канал для статусов - либо структуру, где в заголовочнике
		// будет статус, а в теле данные, но это вроде усложнение
		r := bufio.NewReader(stdout)
		buf := make([]byte, 0, 4<<20)
		for {
			n, err := r.Read(buf[:cap(buf)])
			buf = buf[:n]
			if err != nil {
				if err == io.EOF {
					break // дошли до конца пайпа
				} else {
					chanerr <- err // на проблемку напали
					return
				}
			}
			contentChan <- buf
		}
		close(contentChan)
		close(errorChan)
	}(contentChan, errorChan)

	// дефер использует значения, которые
	// он вычислил на момент вызова дефер,
	// а не на момент вызова ретурн, поэтому
	// дефер вейт только после того, как мы
	// добавили всеъ в вейт груп (????)
	defer func() {
		log.Println("waiting for goroutines...")
		wg.Wait()
	}()

	// основная горутина читает из каналов и следит, чтобы все ок было
	content := make([]byte, 0)
	done := false

	// в чем мем: горутина с чтецом закрывает канал, когда закончила
	// а как нам об этом узнать, чтобы завершить цикл опроса?
	// c помощью оков - если не ок, значит с каналом что-то не так
	// значит, занулляем его (а зачем занулять? потому что закрытый канал
	// ВСЕГДА возвращает занчение при попытке чтения - дефолтное для типа канала)
	// для нас это значит, что если канал с данными закрылся, то в селекте его ветка
	// ВСЕГДА отрабатывает с nil и в другие ветки мы не попадем - закрытость канала
	// можно проверить оком
	for !done {
		select {
		case err, ok := <-errorChan:
			if !ok {
				// если канал закрылся, занулим его
				errorChan = nil
			} else if err != nil {
				// если ты упал, у тебя провал, ты решил, что шансов больше нет
				return nil, err
			}
			// дождаться завершения горутин и выйти
		case chunk, ok := <-contentChan:
			// добавить чанк к нашим чанкам
			if !ok {
				// если канал закрылся, занулим его и завершим цикл опроса
				contentChan = nil
				done = true
			} else {
				content = append(content, chunk...)
			}
		}
	}

	if err := bash.Wait(); err != nil {
		log.Printf("Error waiting for command execution: %s......\n", err.Error())
		return nil, err
	}

	log.Printf("%+v\n", *NewCommandHeader(AddGenome, int64(len(content))))

	// возвращаем в обработчик, чтобы передали бэку
	return &AddGenomeCommand{
		Header: *NewCommandHeader(AddGenome, int64(len(content))),
		Data:   content,
	}, nil
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
