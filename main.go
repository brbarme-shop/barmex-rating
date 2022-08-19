package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/brbarme-shop/brbarmex-rating/server"
)

func main() {

	go func() {
		log.Println(http.ListenAndServe(":9000", nil))
	}()

	server.Start()
}
