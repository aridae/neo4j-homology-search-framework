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

func (cmd *AddGenomeCommand) MarshalBody() ([]byte, int64) {
	return cmd.Data, int64(len(cmd.Data))
}

func (cmd *AddGenomeCommand) UnmarshalBody(b []byte) error {
	cmd.Data = b
	return nil
}
