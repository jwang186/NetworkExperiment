package main

import (
    "fmt"
    "net"
    "net/http"
    "context"
    "io/ioutil"
)

func main() {
    httpClient := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", "/tmp/demo_server.sock")
			},
		},
	}
    response, _ := httpClient.Get("http://anyString/ok")
    data, _ := ioutil.ReadAll(response.Body)
    fmt.Println(string(data))
}
