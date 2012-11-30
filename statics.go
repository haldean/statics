package main

import (
	"flag"
	"net/http"
	"strings"
)

var baseDirectory *string = flag.String("base", ".", "Base directory for files.")
var port *int = flag.Int("port", 8080, "Port to listen for requests on.")

func hostDirectory(r *http.Request) http.Dir {
	hostName := strings.Split(r.Host, ":")[0]
	path := strings.Join(
		[]string{*baseDirectory, strings.Replace(hostName, "/", "", -1)}, "/")
	return http.Dir(path)
}

func hostDispatcher(w http.ResponseWriter, r *http.Request) {
	handler := http.FileServer(hostDirectory(r))
	handler.ServeHTTP(w, r)
}

func main() {
	flag.Parse()
	http.HandleFunc("/", hostDispatcher)
	http.ListenAndServe(":8080", nil)
}
