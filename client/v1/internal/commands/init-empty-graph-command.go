package commands

import (
	"encoding/binary"
	"unsafe"
)

type InitEmptyGraphCommand struct {
	Header CommandHeader
	K      int64
}

func NewInitEmptyGraphCommand(k int64) Command {
	return &InitEmptyGraphCommand{
		K:      k,
		Header: *NewCommandHeader(InitEmptyGraph, int64(unsafe.Sizeof(k))),
	}
}

func (cmd *InitEmptyGraphCommand) GetCmd() CommandOption {
	return cmd.Header.Cmd
}

func (cmd *InitEmptyGraphCommand) GetHeader() CommandHeader {
	return cmd.Header
}

func (cmd *InitEmptyGraphCommand) MarshalBody() ([]byte, int64) {
	bytes := make([]byte, 8)
	binary.PutVarint(bytes, cmd.K)
	return bytes, int64(len(bytes))
}

func (cmd *InitEmptyGraphCommand) UnmarshalBody(b []byte) error {
	cmd.K, _ = binary.Varint(b)
	return nil
}
