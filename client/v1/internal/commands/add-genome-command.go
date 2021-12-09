package commands

import (
	"fmt"
	"os/exec"
	"strconv"
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
	data, err := bash.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// читаем жсон из вывода
	println(string(data), len(data))

	// возвращаем в обработчик, чтобы передали бэку
	return &AddGenomeCommand{
		Header: *NewCommandHeader(AddGenome, int64(len(data))),
		Data:   data,
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
