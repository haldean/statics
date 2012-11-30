package main

import (
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"syscall"
)

var baseDirectory *string = flag.String("base", ".", "Base directory for files.")
var port *int = flag.Int("port", 8080, "Port to listen for requests on.")

var handlerCache map[http.Dir]http.Handler

func Abs(name string) (string, error) {
	if path.IsAbs(name) {
		return name, nil
	}
	wd, err := os.Getwd()
	return path.Join(wd, name), err
}

func hostDirectory(r *http.Request) (http.Dir, error) {
	hostName := strings.Split(r.Host, ":")[0]
	if hostName[0] == '.' {
		return http.Dir(""), errors.New("Illegal host name")
	}

	path := strings.Replace(hostName, "/", "", -1)
	return http.Dir(path), nil
}

func hostDispatcher(w http.ResponseWriter, r *http.Request) {
	dir, err := hostDirectory(r)
	if err != nil {
		http.NotFoundHandler().ServeHTTP(w, r)
		return
	}

	log.Printf("Request for %v %v\n", r.Host, r.URL)
	handler := handlerCache[dir]
	if handler == nil {
		handler = http.FileServer(dir)
		handlerCache[dir] = handler
	}

	handler.ServeHTTP(w, r)
}

func main() {
	flag.Parse()

	absPath, err := Abs(*baseDirectory)
	if err != nil {
		log.Fatalf("Unable to resolve path %v\n", *baseDirectory)
	}

	os.Chdir(absPath)
	log.Printf("Chrooting to %v\n", absPath)
	if err := syscall.Chroot(absPath); err != nil {
		log.Fatalf("Unable to chroot: %v\n", err)
	}

	handlerCache = make(map[http.Dir]http.Handler)
	http.HandleFunc("/", hostDispatcher)
	http.ListenAndServe(":8080", nil)
}
