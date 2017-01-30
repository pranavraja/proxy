package proxy

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"

	"github.com/elazarl/goproxy"
)

func isPrintable(contentType string) bool {
	switch {
	case strings.HasPrefix(contentType, "video/"), strings.HasPrefix(contentType, "audio/"):
		return false
	case strings.HasPrefix(contentType, "image/"):
		return false
	case contentType == "application/mp4":
		return false
	case contentType == "application/octet-stream":
		return false
	default:
		return true
	}
}

func New(processRequest func(*http.Request), processResponse func(*http.Response)) *goproxy.ProxyHttpServer {
	proxy := goproxy.NewProxyHttpServer()
	if *MitmOnly != "" {
		proxy.OnRequest().HandleConnect(goproxy.FuncHttpsHandler(func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
			// e.g. -mitmonly '!itunes.apple.com' mitms all hosts except that one
			if strings.HasPrefix(*MitmOnly, "!") {
				if strings.Contains(strings.TrimPrefix(*MitmOnly, "!"), strings.TrimSuffix(host, ":443")) {
					return goproxy.OkConnect, host
				}
				return goproxy.MitmConnect, host
			}
			if !strings.Contains(*MitmOnly, strings.TrimSuffix(host, ":443")) {
				return goproxy.OkConnect, host
			}
			return goproxy.MitmConnect, host
		}))
	} else if !*NoMITM {
		proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	}
	proxy.OnRequest().DoFunc(func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		if processRequest != nil {
			processRequest(r)
		}

		ctx.UserData = time.Now()
		if OutFile != nil && strings.HasSuffix(r.Host, *Only) {
			d, _ := httputil.DumpRequest(r, isPrintable(r.Header.Get("Content-Type")))
			OutFile.Write(d)
			OutFile.Write([]byte("\n\n"))
		}
		return r, nil
	})

	proxy.OnResponse().DoFunc(func(r *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		if r == nil {
			return r
		}
		if processResponse != nil {
			processResponse(r)
		}
		if OutFile == nil && !strings.HasSuffix(r.Request.Host, *Only) {
			return r
		}
		if !*Quiet {
			if ctx.UserData != nil {
				elapsed := time.Since(ctx.UserData.(time.Time))
				if elapsed > 8*time.Second {
					log.Printf("%s %s %s => %s", ColoredMethod(r.Request.Method), ColoredURL(r.Request.URL), Colored(time.Since(ctx.UserData.(time.Time)).String(), Red), ColoredStatusLine(r.StatusCode, http.StatusText(r.StatusCode)))
				} else {
					log.Printf("%s %s %s => %s", ColoredMethod(r.Request.Method), ColoredURL(r.Request.URL), Colored(time.Since(ctx.UserData.(time.Time)).String(), Cyan), ColoredStatusLine(r.StatusCode, http.StatusText(r.StatusCode)))
				}
			} else {
				log.Printf("%s %s => %s", ColoredMethod(r.Request.Method), ColoredURL(r.Request.URL), ColoredStatusLine(r.StatusCode, http.StatusText(r.StatusCode)))
			}
		}
		if OutFile != nil && strings.HasSuffix(r.Request.Host, *Only) {
			d, _ := httputil.DumpResponse(r, isPrintable(r.Header.Get("Content-Type")))
			fmt.Fprintf(OutFile, "%s %s %s\n", r.Request.Method, r.Request.URL, time.Since(ctx.UserData.(time.Time)))
			OutFile.Write(d)
			OutFile.Write([]byte("\n\n"))
		}
		return r
	})
	return proxy
}

var Addr = flag.String("addr", ":8080", "proxy listen address")
var CertFile = flag.String("cert", os.Getenv("GOPATH")+"/src/github.com/pranavraja/proxy/cert", "cert name e.g. if you have cert.crt and cert.key on disk, the cert is called 'cert'")
var Dump = flag.String("dump", "", "dump all requests to a file")
var Only = flag.String("only", "", "dump only requests matching a hostname suffix")
var NoMITM = flag.Bool("nomitm", false, "don't mitm HTTPS traffic")
var MitmOnly = flag.String("mitmonly", "", "only mitm a certain host")
var Quiet = flag.Bool("quiet", false, "Quiet output")
var OutFile *os.File

func init() {
	flag.Parse()

	var err error
	if *Dump != "" {
		OutFile, err = os.Create(*Dump)
		if err != nil {
			log.Fatal(err)
		}
	}

	cert, err := ioutil.ReadFile(*CertFile + ".crt")
	if err != nil {
		log.Fatal(err)
	}
	key, err := ioutil.ReadFile(*CertFile + ".key")
	if err != nil {
		log.Fatal(err)
	}
	goproxy.GoproxyCa, err = tls.X509KeyPair(cert, key)
	if err != nil {
		log.Fatal(err)
	}
}
