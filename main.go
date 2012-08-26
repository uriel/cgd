package main

import (
	"fmt"
	"os"
	"log"
	"flag"
	"net"
	"net/http"
	"net/http/cgi"
	"net/http/fcgi"
)

var cmd = flag.String("c", "", "CGI program to run")
var pwd = flag.String("w", "", "Working dir for CGI")
var serveFcgi = flag.Bool("f", false, "Run as a FCGI 'server' instead of HTTP")
var address = flag.String("a", ":3333", "Listen address")

func handleCgi(res http.ResponseWriter, req *http.Request) {
	c := *cmd
	if c[0] != "/"[0] {
		c = "./"+c
	}

	os.Setenv("PATH", os.Getenv("PATH")+":.")

	h := cgi.Handler{
		Path: c,
		Root: "/", 
		Dir: *pwd,
		InheritEnv: []string{"PATH", "PLAN9",},
		}
	fmt.Fprintf(os.Stderr, "%v", h)
	h.ServeHTTP(res, req)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if *cmd == "" {
		usage()
	}

	h := http.HandlerFunc(handleCgi)

	if *serveFcgi {
		l, err := net.Listen("tcp", *address)
		if err != nil {
			log.Fatal(err)
		}
		fcgi.Serve(l, h)
	} else {
		http.Handle("/", h)
		http.ListenAndServe(*address, nil)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: cgd [-s] -c prog [-w wdir] [-a addr]")

	flag.PrintDefaults()
	os.Exit(2)
}
