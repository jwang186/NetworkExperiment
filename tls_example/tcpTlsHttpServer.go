package main

import (
	"net/http"
	"fmt"
)

func httpHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	body := "\033[92m[Success]Reached TCP Server" + "\n"
	w.Write([]byte(body))
}

func main() {
	http.HandleFunc("/ok", httpHandler)
	fmt.Println("Server started, supported URL path: /ok")
	err := http.ListenAndServeTLS(":443", "/tmp/server-chain.crt", "/tmp/server.key", nil)
	if err != nil {
      fmt.Print("ListenAndServe: ", err)
    }
}
