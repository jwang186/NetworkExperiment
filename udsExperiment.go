package main

import (
	"net/http"
	"net"
	"fmt"
    "context"
    "io/ioutil"
    "syscall"
    "os"
    "time"
    "sync"
    "os/signal"

	 retryablehttp "github.com/hashicorp/go-retryablehttp"
)

const(
    envoyAminUdsPath = "/tmp/envoy_admin.sock"
    appUdsPath = "/tmp/app_admin.sock"
    envoyFile = "/tmp/envoy-uds-admin.yaml"
)

type MessageSource struct {
    pid int
}

func (ms *MessageSource) GetPid() int {
    return ms.pid
}

func (ms *MessageSource) SetPid(id int) {
    ms.pid = id
}

type EnvoyEndpointHttpData struct {
    client *retryablehttp.Client
}

func setupUdsPath() error {
    _, err := os.Create(envoyAminUdsPath)
    if err != nil && !os.IsExist(err){
        fmt.Println("Failed to create uds file %s", envoyAminUdsPath)
        return err
    }
    return nil
}

func buildHttpClient() *retryablehttp.Client {
    httpClient := http.Client{
    	Transport: &http.Transport{
    		DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
    			return net.Dial("unix", envoyAminUdsPath)
    		},
    	},
    }
    retryClient := &retryablehttp.Client{
        Logger: nil,
        HTTPClient:   &httpClient,
        RetryWaitMin: 1,
        RetryWaitMax: 5,
        RetryMax:     8,
        CheckRetry:   retryablehttp.DefaultRetryPolicy,
        Backoff:      retryablehttp.DefaultBackoff,
    }
    return retryClient
}

func (envoyEndpoint *EnvoyEndpointHttpData) StatsHandler(w http.ResponseWriter, r *http.Request) {
	request, _ := retryablehttp.NewRequest("GET", "http://127.0.0.1:9901/stats", nil)
    request.Header.Add("Connection", "close")
    request.Header.Add("User-Agent", "AppNet")
    response, err := envoyEndpoint.client.Do(request)
    if err != nil {
        fmt.Println(err)
        return
    }
    responseBody, err := ioutil.ReadAll(response.Body)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(string(responseBody))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(string(responseBody)))
}

func serverOkHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("The HTTP Server is working fine!"))
}

func startHttpServer() {
    os.Remove(appUdsPath)
    listener, err := net.Listen("unix", appUdsPath)
    if err != nil {
    	fmt.Println(err)
    }
    defer listener.Close()
    defer os.Remove(appUdsPath)
    var envoyEndpoint EnvoyEndpointHttpData
    envoyEndpoint.client = buildHttpClient()

    http.HandleFunc("/stats", envoyEndpoint.StatsHandler)
    http.HandleFunc("/ok", serverOkHandler)
    http.Serve(listener, nil)
}

func startEnvoy(ms *MessageSource) {
    attr := syscall.ProcAttr{
        Dir: ".",
        Env: os.Environ(),
        Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()},
        Sys: nil,
    }
    pid, err := syscall.ForkExec("/usr/bin/envoy", []string{"/usr/bin/envoy", "-c", envoyFile}, &attr)
    if err != nil || pid == -1 || pid == 0 {
        fmt.Println(fmt.Sprintf("Unable to start envoy %v", err))
    }
    ms.SetPid(pid)
    dir, _ := syscall.Getwd()
    fmt.Println(fmt.Sprintf("Envoy started. Working dir [%s]. Pid: [%d]", dir, pid))
}

func stopEnvoy(ms *MessageSource) {
    maxWaitTime := time.Duration(20*time.Second)
    var wg sync.WaitGroup
    pid := ms.GetPid()

    wg.Add(1)
    go func() {
        var processActive bool
        var startTime = time.Now()
        	// Check process state every second
        	ticker := time.NewTicker(1 * time.Second)
        	defer ticker.Stop()
        	defer wg.Done()

        	for processActive {
        		select {
        		case <-ticker.C:
        		    _, err := os.FindProcess(int(pid))
        		    if err != nil {
        		        processActive = false
        		    }
        			if !processActive {
        				return
        			}

        			if time.Since(startTime) > maxWaitTime {
        				processActive = false
        			}
        		}
        	}

        	pid = ms.GetPid()
        	if pid > 0 {
        		fmt.Println("Killing pid ", pid)
        		syscall.Kill(pid, syscall.SIGKILL)
        		ms.SetPid(-1)
        		processActive = false
        	}
    }()

    wg.Wait()
}

func main() {
    var ms MessageSource

    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        sig := <-sigs
        fmt.Println()
        fmt.Println(sig)
        pid := ms.GetPid()
        fmt.Println("Killing pid ", pid)
        //syscall.Kill(pid, syscall.SIGTERM)
    }()

    //go startEnvoy(&ms)
    //defer stopEnvoy(&ms)

	startHttpServer()
}