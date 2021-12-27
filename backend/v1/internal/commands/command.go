package commands

import (
	"encoding/binary"
	"encoding/json"
	"log"
)

const (
	InitEmptyGraph CommandOption = iota
	AddGenome      CommandOption = iota
)

var (
	HeaderLen = 47
)

type Command interface {
	GetCmd() CommandOption
	GetHeader() CommandHeader
	MarshalBody() ([]byte, int64)
	UnmarshalBody([]byte) error
}

type CommandOption int64

type CommandHeader struct {
	Cmd      CommandOption `json:"cmd"`
	BodySize int64         `json:"bodylen"`
}

type CommandHeaderDTO struct {
	Cmd      []byte `json:"cmd"`
	BodySize []byte `json:"bodylen"`
}

// го маршаллит int64 как строку, а надо - как 8 байт,
// чтобы у хедера была фиксированная длина
func (h CommandHeader) Marshall() ([]byte, error) {
	bodylenBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bodylenBytes, uint64(h.BodySize))

	cmdBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(cmdBytes, uint64(h.Cmd))

	m := CommandHeaderDTO{
		Cmd:      cmdBytes,
		BodySize: bodylenBytes,
	}

	bytes, err := json.Marshal(m)
	if err != nil {
		log.Printf("failed to marshall header %+v - %s\n", h, err)
		return nil, err
	}

	log.Printf("Marshalling header result: %d\n", len(bytes))
	return bytes, nil
}

func UnmarshallHeader(bytes []byte) (*CommandHeader, error) {

	m := CommandHeaderDTO{}
	err := json.Unmarshal(bytes, &m)
	if err != nil {
		log.Printf("*** failed to unmarshall header %s\n", err)
		return nil, err
	}

	bodylen := int64(binary.LittleEndian.Uint64(m.BodySize))
	cmd := int64(binary.LittleEndian.Uint64(m.Cmd))

	header := &CommandHeader{
		Cmd:      CommandOption(cmd),
		BodySize: bodylen,
	}
	log.Printf("Unarshalling header result: %+v\n", header)
	return header, nil
}

func NewCommandHeader(cmd CommandOption, bodySize int64) *CommandHeader {
	return &CommandHeader{
		Cmd:      cmd,
		BodySize: bodySize,
	}
}
