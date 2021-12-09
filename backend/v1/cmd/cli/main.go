package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	cmd "github.com/aridae/neo4j-homology-search-framework/backend/v1/internal/commands"
	cnt "github.com/aridae/neo4j-homology-search-framework/backend/v1/internal/controller"
	dataaccess "github.com/aridae/neo4j-homology-search-framework/backend/v1/internal/data-access"
	db "github.com/aridae/neo4j-homology-search-framework/backend/v1/internal/dbdriver"
	proto "github.com/aridae/neo4j-homology-search-framework/backend/v1/internal/proto"
	"github.com/joho/godotenv"
)

const (
	ListenerToControllerBuf = 3
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("no .env file found, using default env variables")
	}
}

func main() {

	ctx, cancelFunc := context.WithCancel(context.Background())

	// ожидаем коллекцию горутин из мейна:
	// тут контроллер и листенер
	wg := &sync.WaitGroup{}

	// создаем клиента для бд
	neo4jClient, err := db.GetNeo4jClient(db.GetNeo4jOptions())
	if err != nil {
		log.Println(err)
		return
	}

	// перед выходом почистим
	defer neo4jClient.CloseNeo4jClient()

	// создаем канал для общения листенера и контроллера
	commandsChan := make(chan cmd.Command, ListenerToControllerBuf)

	// вешаем на сокет листенера
	listener := proto.OpenListener("localhost", 35035, commandsChan)

	// заводим контроллера
	controller := cnt.NewController(dataaccess.NewNeo4jRepository(neo4jClient), commandsChan)

	// ловим сигналы
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
	wg.Add(1)
	go func(context context.Context) {
		defer wg.Done()

		// горутина будет заблокирована на чтении
		// до тех пор, пока мы не пошлем ей сигнал
		s := <-sigs
		log.Printf("RECEIVED SIGNAL [SIGTERM]: %s", s)

		// передаем всем запущенным горутинам завершиться
		cancelFunc()
	}(ctx)

	// начать работу контроллера
	wg.Add(1)
	go func(context context.Context) {
		defer wg.Done()
		controller.RunControl(context)
	}(ctx)

	// начать работу хэндлера
	wg.Add(1)
	go func(context context.Context) {
		defer wg.Done()
		listener.Serve(context)
	}(ctx)

	wg.Wait()
	log.Println("Bye!")

}
