package main

import (
	"log"
	"net/http"
	"os"

	"github.com/pranavraja/proxy"
)

func main() {
	if len(os.Args) <= 1 || len(os.Args)%2 == 0 {
		log.Fatal("usage: map host1 other host2 other")
	}
	mappings := make(map[string]string)
	for i := 1; i < len(os.Args); i += 2 {
		mappings[os.Args[i]] = os.Args[i+1]
	}
	p := proxy.New(func(r *http.Request) {
		for host, remote := range mappings {
			if r.Host == host {
				r.Host = remote
				r.URL.Host = remote
				r.URL.Scheme = "http"
			}
		}
	}, nil)
	log.Fatal(http.ListenAndServe(*proxy.Addr, p))
}
