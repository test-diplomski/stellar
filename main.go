package main

import (
	"github.com/c12s/stellar/model"
	"github.com/c12s/stellar/server"
	"github.com/c12s/stellar/storage/etcd"
	"log"
	"time"
)

func main() {

	// Load configurations
	conf, err := model.ConfigFile()
	if err != nil {
		log.Fatal(err)
	}

	//Load database
	db, err := etcd.New(conf, 10*time.Second)
	if err != nil {
		log.Fatal(err)
	}

	//Start Server
	server.Run(conf, db)
}
