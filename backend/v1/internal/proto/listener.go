package listener

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"

	cmd "github.com/aridae/neo4j-homology-search-framework/backend/v1/internal/commands"
)

type Listener struct {
	Host string
	Port int64

	// мы только слушаем и парсим,
	// а на выполнение отдаем во внешний канал
	listener *net.Listener

	// после того, как хэндлер получил команду
	// он приводит ее в божеский вид, парсит и отдает контроллеру
	// контроллер поделит ее на таски и раздаст воркерам из пула
	// специально не вложила контроллер в листенера,
	// потому что там терки с многопоточностью получаются -
	// у листенера свои потоки, у контроллера свои
	// в мультиплексорах это как-то реализовано нормально, мб посмотреть(!!!!)
	OutQueue chan cmd.Command
}

func OpenListener(Host string, Port int64, Out chan cmd.Command) *Listener {
	listener, err := net.Listen("tcp", Host+":"+strconv.FormatInt(Port, 10))
	if err != nil {
		log.Fatal("failed to init listener")
	}

	//log.Printf("Listening on %s", listener.Addr().String())
	return &Listener{
		Host:     Host,
		Port:     Port,
		listener: &listener,
		OutQueue: Out,
	}
}

func (l *Listener) IsOpen() bool {
	return l.listener != nil
}

func (l *Listener) Close() {
	if l.IsOpen() {
		(*l.listener).Close()
	}
}

// serving loop
func (l *Listener) Serve(ctx context.Context) {
	log.Printf("Serving port %s", (*l.listener).Addr().String())
	connsChan := make(chan net.Conn, 10)

	defer log.Println("Listener serving loop stopped")
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(listener net.Listener) {
		defer wg.Done()
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Println("Error accepting: ", err.Error())
				return
			}
			log.Printf("Accepted new connection: %s\n", conn.LocalAddr().String())
			connsChan <- conn
		}
	}(*l.listener)

	for {
		select {
		case conn := <-connsChan:
			// обрабатываем новое соединение (в той же!!! горутине)
			// потому действие короткое и на оверхед больше потратится
			log.Printf("New connection: %s\n", conn.LocalAddr().String())
			l.handleCommand(conn)
		case <-ctx.Done():
			log.Println("closing listener...")
			(*l.listener).Close()
			wg.Wait()
			return
		}
	}
}

func (l *Listener) readCommand(conn net.Conn) (cmd.Command, error) {
	// read command header - parse, get type
	headerBytes := make([]byte, cmd.HeaderLen)
	_, err := conn.Read(headerBytes)
	if err != nil {
		log.Printf("Failed to read header from connection %s\n", conn.LocalAddr().String())
		return nil, err
	}

	// unserialize it
	var Header cmd.CommandHeader
	err = json.Unmarshal(headerBytes, &Header)
	if err != nil {
		log.Println("Failed to unmarshal header:", err)
		return nil, err
	}
	log.Printf("read from socket - %+v (%d)", Header, len(headerBytes))

	// resolve command type
	var command cmd.Command
	switch Header.Cmd {
	case cmd.InitEmptyGraph:
		command = &cmd.InitEmptyGraphCommand{
			Header: Header,
		}
	case cmd.AddGenome:
		command = &cmd.AddGenomeCommand{
			Header: Header,
		}
	default:
		return nil, fmt.Errorf("failed to resolve command: unsopported command")
	}

	// read command body
	bodyBytes := make([]byte, Header.BodySize)
	_, err = conn.Read(bodyBytes)
	if err != nil {
		log.Printf("Failed to read body from connection %s\n", conn.LocalAddr().String())
		return nil, err
	}

	// unserialize it
	err = command.UnmarshalBody(bodyBytes)
	if err != nil {
		log.Printf("Failed to read body from connection %s\n", conn.LocalAddr().String())
		return nil, err
	}
	return command, nil
}

func (l *Listener) handleCommand(conn net.Conn) error {
	command, err := l.readCommand(conn)
	if err != nil {
		log.Println("failed to read command from socket", command)
		return err
	}

	// отправляем контроллеру,
	// в го каналы конкуррентобезопасны
	// так что проблем быть не должно
	l.OutQueue <- command
	return nil
}
