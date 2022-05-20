package main

import (
	"net/http"
	"net"
	"fmt"
)

func udsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	body := "Reached UDS Server" + "\n"
	w.Write([]byte(body))
}

func main() {
	listener, _ := net.Listen("unix", "demo.sock")
        defer listener.Close()
	http.HandleFunc("/uds", udsHandler)
	fmt.Println("Server started, supported URL path: /uds")
	http.Serve(listener, nil)
}
