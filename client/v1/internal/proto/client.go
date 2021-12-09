package proto

import (
	"encoding/json"
	"log"
	"net"

	cmd "github.com/aridae/neo4j-homology-search-framework/client/v1/internal/commands"
)

type Client struct {
	ServerAddr string
}

func NewClient(serverAddr string) *Client {
	return &Client{
		ServerAddr: serverAddr,
	}
}

func (client *Client) SendInitDBCommand(k int64) error {
	initdbCommand := cmd.NewInitEmptyGraphCommand(k)
	return client.sendCommand(initdbCommand)
}

func (client *Client) SendAddGenomeCommand(path string) error {
	addGenomeCommand := cmd.NewAddGenomeCommand(path)
	return client.sendCommand(addGenomeCommand)
}

func (client *Client) sendCommand(cmd cmd.Command) error {
	log.Printf("Connecting to %s...\n", client.ServerAddr)
	conn, err := net.Dial("tcp", client.ServerAddr)
	if err != nil {
		log.Printf("failed to connect to %s: %s", client.ServerAddr, err)
		return err
	}

	// отправляем хедер
	log.Printf("trying to send header %+v\n", cmd.GetHeader())
	CommandHeader, err := json.Marshal(cmd.GetHeader())
	if err != nil {
		log.Println("failed to marshall command: ", err)
		return err
	}
	_, err = conn.Write(CommandHeader)
	if err != nil {
		log.Println("failed to send header: ", err)
		return err
	}

	// отправляем боди
	commandBody, bodyLen := cmd.MarshalBody()
	if bodyLen == 0 {
		log.Println("failed to marshal body: ", err)
		return err
	}
	log.Printf("trying to send body %d\n", bodyLen)
	_, err = conn.Write(commandBody)
	if err != nil {
		log.Println("failed to send body: ", err)
		return err
	}

	conn.Close()
	return nil
}
