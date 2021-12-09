package commands

import (
	"encoding/json"
	"log"
)

type AddGenomeCommand struct {
	Header CommandHeader
	Path   string
}

func NewAddGenomeCommand(path string) Command {
	return &AddGenomeCommand{
		Path:   path,
		Header: *NewCommandHeader(AddGenome, int64(len(path))),
	}
}

func (cmd *AddGenomeCommand) GetCmd() CommandOption {
	return cmd.Header.Cmd
}

func (cmd *AddGenomeCommand) GetHeader() CommandHeader {
	return cmd.Header
}

// func (cmd *AddGenomeCommand) MarshalHeader() ([]byte, int64) {
// 	bytes := []byte(cmd.Path)
// 	return bytes, int64(len(bytes))
// }

// func (cmd *AddGenomeCommand) UnmarshalHeader(b []byte) error {
// 	err := json.Unmarshal(b, &cmd.Path)
// 	if err != nil {
// 		log.Println("failed to unmarshall body:", err)
// 		return err
// 	}
// 	return nil
// }

func (cmd *AddGenomeCommand) MarshalBody() ([]byte, int64) {
	bytes := []byte(cmd.Path)
	return bytes, int64(len(bytes))
}

func (cmd *AddGenomeCommand) UnmarshalBody(b []byte) error {
	err := json.Unmarshal(b, &cmd.Path)
	if err != nil {
		log.Println("failed to unmarshall body:", err)
		return err
	}
	return nil
}
