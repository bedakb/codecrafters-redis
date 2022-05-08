package main

import (
	"log"

	"github.com/bedakb/codecrafters-redis/server"
)

func main() {
	port := ":6379"
	log.Printf("running redis-sever on port %s", port)
	log.Fatal(server.ListenAndServe(port))
}
