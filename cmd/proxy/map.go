package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/pranavraja/proxy"
)

type hostpath struct {
	Host string
	Path string
}

func parseHostpath(s string) hostpath {
	hp := strings.SplitN(s, "/", 2)
	m := hostpath{Host: hp[0]}
	if len(hp) > 1 {
		m.Path = "/" + hp[1]
	}
	return m
}

func main() {
	if len(os.Args) <= 1 || len(os.Args)%2 == 0 {
		log.Fatal("usage: map host1 other host2 other")
	}
	mappings := make(map[hostpath]hostpath)
	for i := 1; i < len(os.Args); i += 2 {
		mappings[parseHostpath(os.Args[i])] = parseHostpath(os.Args[i+1])
	}
	p := proxy.New(func(r *http.Request) {
		for og, remote := range mappings {
			if r.Host == og.Host && strings.HasPrefix(r.URL.Path, og.Path) && !strings.HasPrefix(r.URL.Path, remote.Path) {
				r.Host = remote.Host
				r.URL.Host = remote.Host
				r.URL.Path = remote.Path + strings.TrimPrefix(r.URL.Path, og.Path)
				r.URL.Scheme = "http"
			}
		}
	}, nil)
	log.Fatal(http.ListenAndServe(*proxy.Addr, p))
}
