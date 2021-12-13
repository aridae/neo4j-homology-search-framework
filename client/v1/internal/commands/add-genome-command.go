package commands

import (
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
		io.Copy(os.Stdout, stdout)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		io.Copy(os.Stderr, stderr)
	}()

	wg.Wait()

	if err := bash.Wait(); err != nil {
		log.Printf("Error waiting for command execution: %s......\n", err.Error())
		return nil, err
	}

	// возвращаем в обработчик, чтобы передали бэку
	return &AddGenomeCommand{
		Header: *NewCommandHeader(AddGenome, int64(0)),
		Data:   []byte{},
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
