// A stand-alone HTTP server providing a web UI for node management.
package main

import (
	"net/http"

	"github.com/tracyde/hiera-gui/server"
)

func main() {
	server.RegisterHandlers()
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.ListenAndServe(":9080", nil)
}
