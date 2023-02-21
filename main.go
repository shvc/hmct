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
	"github.com/gorilla/mux"
)

var (
	debug      bool
	version    string = "0.0"
	addr       string = ":80"
	msg        string = "default message"
	dataDir    string = os.TempDir()
	configFile string = "config.json"
)

func logRequest(r *http.Request, status int) {
	if debug {
		log.Println(r.RemoteAddr, r.Method, r.RequestURI, r.ContentLength, r.Host, status)
		return
	}
	log.Println(r.RemoteAddr, r.Method, r.RequestURI, r.ContentLength, status)
}

func router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/upload/{name:.+}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		filename := filepath.Join(dataDir, vars["name"])
		fd, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
		if err != nil {
			logRequest(r, http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
		defer fd.Close()
		io.Copy(fd, r.Body)
		logRequest(r, http.StatusOK)
	}).Methods(http.MethodPut)

	router.HandleFunc("/download/{name:.+}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		filename := filepath.Join(dataDir, vars["name"])
		fd, err := os.Open(filename)
		if err != nil {
			logRequest(r, http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
		defer fd.Close()
		io.Copy(w, fd)
		logRequest(r, http.StatusOK)
	}).Methods(http.MethodGet)

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		logRequest(r, http.StatusOK)
		w.Write([]byte("ok"))
	})

	router.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		if hn, err := os.Hostname(); err == nil {
			w.Write([]byte(fmt.Sprintf("hostname: %v\n", hn)))
		}
		w.Write([]byte(fmt.Sprintf("version: %v\n", version)))
		w.Write([]byte(fmt.Sprintf("config: %v\n", configFile)))
		w.Write([]byte(fmt.Sprintf("msg: %v\n", msg)))
		logRequest(r, http.StatusOK)
	})

	router.HandleFunc("/env", func(w http.ResponseWriter, r *http.Request) {
		logRequest(r, http.StatusOK)
		for _, es := range os.Environ() {
			w.Write([]byte(fmt.Sprintln(es)))
		}
	})

	router.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
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
	return router
}

func main() {
	ver := flag.Bool("v", false, "show version")
	flag.BoolVar(&debug, "debug", debug, "debug log level")
	flag.StringVar(&msg, "msg", msg, "server message")
	flag.StringVar(&addr, "addr", addr, "server serve address")
	flag.StringVar(&dataDir, "data-dir", dataDir, "server data dir")
	flag.StringVar(&configFile, "config", configFile, "server config file")
	envflag.Parse()

	if *ver {
		fmt.Println("version", version)
		return
	}

	log.Println(filepath.Base(os.Args[0]), version, "listen and serve", addr)
	if err := http.ListenAndServe(addr, router()); err != nil {
		log.Println("listen and serve error", err)
	}
}
