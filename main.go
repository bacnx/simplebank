package main

import (
	"database/sql"
	"log"

	"github.com/bacnx/simplebank/api"
	db "github.com/bacnx/simplebank/db/sqlc"
	"github.com/bacnx/simplebank/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.GetConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(config, store)

	server.Start(config.ServerAddress)
}
