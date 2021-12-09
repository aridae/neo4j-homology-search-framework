package commands

type AddGenomeCommand struct {
	Header CommandHeader
	Data   []byte
}

func NewAddGenomeCommand(data []byte) Command {
	return &AddGenomeCommand{
		Data:   data,
		Header: *NewCommandHeader(AddGenome, int64(len(data))),
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
	return cmd.Data, int64(len(cmd.Data))
}

func (cmd *AddGenomeCommand) UnmarshalBody(b []byte) error {
	cmd.Data = b
	return nil
}
