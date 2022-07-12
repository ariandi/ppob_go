package main

import (
	"database/sql"
	"github.com/ariandi/ppob_go/api"
	db "github.com/ariandi/ppob_go/db/sqlc"
	"github.com/ariandi/ppob_go/util"
	_ "github.com/lib/pq"
	"log"
)

const dbDriver = "postgres"
const dbSource = "postgresql://postgres:postgres@localhost:5432/ppob?sslmode=disable"
const serverAddress = "0.0.0.0:8080"

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config : ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db", err)
	}

	store := db.NewStore(conn)
	//server, err := api.NewServer(config, store)
	server := api.NewServer(store)
	//if err != nil {
	//	log.Fatal("cannot create server", err)
	//}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot connect to server", err)
	}
}
