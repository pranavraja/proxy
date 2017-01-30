package main

import (
	"log"
	"net/http"

	"github.com/pranavraja/proxy"
)

func main() {
	p := proxy.New(nil, nil)
	log.Fatal(http.ListenAndServe(*proxy.Addr, p))
}
