package commands

const (
	InitEmptyGraph CommandOption = iota
	AddGenome      CommandOption = iota
)

var (
	HeaderLen = 8 + 8 + 2 + 2
)

type Command interface {
	GetCmd() CommandOption
	GetHeader() CommandHeader
	MarshalBody() ([]byte, int64)
	UnmarshalBody([]byte) error
}

type CommandOption int64

type CommandHeader struct {
	Cmd      CommandOption
	BodySize int64
}

func NewCommandHeader(cmd CommandOption, bodySize int64) *CommandHeader {
	return &CommandHeader{
		Cmd:      cmd,
		BodySize: bodySize,
	}
}
