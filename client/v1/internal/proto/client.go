package proto

import (
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

func (client *Client) SendAddGenomeCommand(path string, k int64) error {
	// дергаем скрипт с кмс, читаем получившийся жсон
	// удаляем этот жсон(?) оставляем в качестве кэша(?)
	// потом удалим...
	_, err := cmd.NewAddGenomeCommand(path, k)
	if err != nil {
		return err
	}

	return nil
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
	CommandHeader, err := cmd.GetHeader().Marshall()
	if err != nil {
		log.Println("failed to marshall header: ", err)
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
	log.Printf("trying to send body %s - %d\n", commandBody[:200], bodyLen)
	_, err = conn.Write(commandBody)
	if err != nil {
		log.Println("failed to send body: ", err)
		return err
	}

	conn.Close()
	return nil
}
