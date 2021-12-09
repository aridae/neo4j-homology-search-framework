package main

import (
	"os"

	"github.com/aridae/neo4j-homology-search-framework/client/v1/internal/proto"
	"gopkg.in/alecthomas/kingpin.v2"
)

// тут маленький ничего не умеющий консольный клиент
// он умеет подключаться к сокету бэка и отправлять команду
// он не умеет даже получать статус выполнения, вот такой он неразумный

// по хорошему на стороне клиента надо завести базу со статусами пендинг реквестов,
// чтобы не держать клиента все время запущенным

var (
	app     = kingpin.New("mewo4j", "A dummy dum-dum command-line miserable neo4j client.")
	connstr = app.Flag("mewoserver", "Server address (connection string).").Required().String()
	initdb  = app.Command("init", "Init empty de bruijn graph database.")
	k       = initdb.Flag("k", "Kmer size.").Required().Int64()

	add  = app.Command("add", "Add genome to database.")
	path = add.Flag("path", "Path to file with genome.").Required().String()
	kAdd = add.Flag("k", "Kmer size.").Required().Int64()
)

func main() {
	var client *proto.Client
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {

	case initdb.FullCommand():
		client = proto.NewClient(*connstr)
		client.SendInitDBCommand(*k)

	case add.FullCommand():
		client = proto.NewClient(*connstr)
		client.SendAddGenomeCommand(*path, *kAdd)
	}
}
