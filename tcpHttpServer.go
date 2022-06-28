package main

import (
	"net/http"
	"net"
	"fmt"
)

func httpHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	body := "\033[92m[Success]Reached TCP Server" + "\n"
	w.Write([]byte(body))
}

func main() {
	listener, _ := net.Listen("tcp", "0.0.0.0:8085")
    defer listener.Close()
	http.HandleFunc("/ok", httpHandler)
	fmt.Println("Server started, supported URL path: /ok")
	http.Serve(listener, nil)
}
