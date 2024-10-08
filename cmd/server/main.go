package main

import (
	"github.com/dsocolobsky/reddys/pkg/database"
	server "github.com/dsocolobsky/reddys/pkg/server"
)

func main() {
	handler := server.NewHandler(database.NewMapDatabase(), database.NewAOF("database.aof"))
	reddysServer := server.NewServer(handler)
	defer reddysServer.Close()
	reddysServer.Serve()
}
