package main

import (
	"net/http"
	"net"
	"fmt"
)

func udsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	body := "\033[92m[Success]Reached UDS Server" + "\n"
	w.Write([]byte(body))
}

func main() {
	listener, _ := net.Listen("unix", "/tmp/demo_server.sock")
        defer listener.Close()
	http.HandleFunc("/ok", udsHandler)
	fmt.Println("Server started, supported URL path: /ok")
	http.Serve(listener, nil)
}
