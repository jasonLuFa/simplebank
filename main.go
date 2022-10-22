package main

import (
	"database/sql"
	"log"

	"github.com/jasonLuFa/simplebank/api"
	db "github.com/jasonLuFa/simplebank/db/sqlc"
	"github.com/jasonLuFa/simplebank/util"

	_ "github.com/lib/pq"
)

// const (
// 	dbDriver = "postgres"
// 	dbSource = "postgresql://admin:admin@localhost:5433/simple_bank?sslmode=disable"
// 	serverAddress = "0.0.0.0:8080"
// )

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:",err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil{
		log.Fatal("cannot start server:",err)
	}


}