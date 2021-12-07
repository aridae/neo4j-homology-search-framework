package main

/*
  -config string
    	path to app config file
  -cmd string
        init-empty - create de Bruijn Graph Template
			-k int - k-mer size, 31 default
		TODO: annotate - add information about given sequence(s) to the database
			-dir string - directory with sequence files
			-file string - single sequence file
		TODO: query - find polymorphisms to the given sequence
			-sequence string
*/

/*
TODO на сегодня:
	* Построить пустой граф де Брюйна в базе
	* Загрузить геномы на сервер
	* Добавить информацию о геномах в базу
*/

import (
	"log"
	"os"
	"os/signal"

	db "github.com/aridae/neo4j-homology-search-framework/backend/v1/internal/dbdriver"
	ctrl "github.com/aridae/neo4j-homology-search-framework/backend/v1/internal/de-bruijn-graph/controller"
	dac "github.com/aridae/neo4j-homology-search-framework/backend/v1/internal/de-bruijn-graph/data-access"
	"github.com/aridae/neo4j-homology-search-framework/backend/v1/internal/utils/workerspool"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("no .env file found")
	}
}

func main() {
	pool := workerspool.GetWorkersPool(8, 10)
	// wg := &sync.WaitGroup{}
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	pool.RunBackground()
	// }()

	// setup signal catching
	sigs := make(chan os.Signal, 1)

	// catch all signals since not explicitly listing
	signal.Notify(sigs)

	neo4jClient, err := db.GetNeo4jClient(db.GetNeo4jOptions())
	if err != nil {
		log.Println(err)
		return
	}
	defer neo4jClient.CloseNeo4jClient()

	// method invoked upon seeing signal
	go func(client *db.Neo4jClient) {
		s := <-sigs
		log.Printf("RECEIVED SIGNAL: %s", s)
		log.Println("Closing neo4j client...")
		client.CloseNeo4jClient()
		os.Exit(1)
	}(neo4jClient)

	repo := dac.NewNeo4jRepository(neo4jClient)
	controller := ctrl.NewController(repo, pool)
	log.Println("Initializing database...")
	controller.RunTask(ctrl.InitEmptyDBG)
	log.Println("That's all, folks!")
}
