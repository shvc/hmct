package main

import (
	"flag"
	"log"
	"os"
	"fmt"
	"path/filepath"
	"net/http"
	"github.com/gobike/envflag"
)

var (
	debug bool
	version string = "0.0"
	addr string = ":80"
	msg string = "default message"
)


func main() {
	ver := flag.Bool("v", false, "show version")
	flag.BoolVar(&debug, "debug", debug, "debug log level")
	flag.StringVar(&msg, "msg", msg, "server message")
	flag.StringVar(&addr, "addr", addr, "server serve address")
	envflag.Parse()

	if *ver {
		fmt.Println("version", version)
		return
	}

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RemoteAddr, r.Method, r.RequestURI)
		w.Write([]byte("pong"))
	})

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RemoteAddr, r.Method, r.RequestURI)
		w.Write([]byte("ok"))
	})

	http.HandleFunc("/msg", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RemoteAddr, r.Method, r.RequestURI)
		w.Write([]byte(msg))
	})

	http.HandleFunc("/ver", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RemoteAddr, r.Method, r.RequestURI)
		w.Write([]byte(version))
	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RemoteAddr, r.Method, r.RequestURI)
		w.Write([]byte(version))
	})

	log.Println(filepath.Base(os.Args[0]), version, "listen and serve", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Println("listen and serve error", err)
	}
}
