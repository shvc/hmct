package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gobike/envflag"
)

var (
	debug      bool
	version    string = "0.0"
	addr       string = ":80"
	msg        string = "default message"
	configFile string = "config.json"
)

func logRequest(r *http.Request, status int) {
	if debug {
		log.Println(r.RemoteAddr, r.Method, r.RequestURI, r.ContentLength, r.Host, status)
		return
	}
	log.Println(r.RemoteAddr, r.Method, r.RequestURI, r.ContentLength, status)
}

func router() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		logRequest(r, http.StatusOK)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		logRequest(r, http.StatusOK)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	http.HandleFunc("/msg", func(w http.ResponseWriter, r *http.Request) {
		logRequest(r, http.StatusOK)
		w.Write([]byte(msg))
	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		logRequest(r, http.StatusOK)
		w.Write([]byte(version))
	})

	http.HandleFunc("/env", func(w http.ResponseWriter, r *http.Request) {
		logRequest(r, http.StatusOK)
		for _, es := range os.Environ() {
			w.Write([]byte(es))
		}
	})

	http.HandleFunc("/config-file", func(w http.ResponseWriter, r *http.Request) {
		logRequest(r, http.StatusOK)
		w.Write([]byte(configFile))
	})

	http.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		fd, err := os.Open(configFile)
		if err != nil {
			logRequest(r, http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		defer fd.Close()
		logRequest(r, http.StatusOK)
		io.Copy(w, fd)
	})
}

func main() {
	ver := flag.Bool("v", false, "show version")
	flag.BoolVar(&debug, "debug", debug, "debug log level")
	flag.StringVar(&msg, "msg", msg, "server message")
	flag.StringVar(&addr, "addr", addr, "server serve address")
	flag.StringVar(&configFile, "config", configFile, "server config file")
	envflag.Parse()

	if *ver {
		fmt.Println("version", version)
		return
	}

	router()
	log.Println(filepath.Base(os.Args[0]), version, "listen and serve", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Println("listen and serve error", err)
	}
}
