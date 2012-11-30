package main

import (
  "flag"
  "net/http"
  "strings"
)

var baseDirectory *string = flag.String("base", ".", "Base directory for files.")
var port *int = flag.Int("port", 8080, "Port to listen for requests on.")

func hostDirectory(r *http.Request) http.Dir {
  return http.Dir(strings.Join(
    []string{ *baseDirectory, strings.Replace(r.Host, "/", "", -1) }, "/"))
}

func hostDispatcher(w http.ResponseWriter, r *http.Request) {
  handler := http.FileServer(hostDirectory(r))
  handler.ServeHTTP(w, r)
}

func main() {
  http.HandleFunc("/", hostDispatcher)
  http.ListenAndServe(":8080", nil)
}
