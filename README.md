# NetworkExperiment
Service connect network experiment with UDS, Container, Envoy and AppNet Agent

## Experiment 1: UDS HttpServer
```go run udsHttpServer.go```

```curl --unix-socket /tmp/demo_server.sock any/ok```

## Experiment 2: UDS HttpClient
In another terminal
```go run udsHttpClient.go``` will give us Reached UDS server


## Experiment 3: Envoy with UDS

```envoy -c demo.yaml```

```curl --unix-socket /tmp/demo_listener.sock any/demo/server/ok```

```curl --unix-socket /tmp/demo_listener.sock any/demo/admin/listeners```

## Experiment 4: Container with UDS
```./start_appnet_agent.sh``` before running this command, make sure we have appnet-agent container image ready and the AWS credentials are set.
This command will start a appnet-agent container and serve a HttpServer on UDS /tmp/appnet_admin.sock

```docker run -v /tmp:/tmp demo``` it will start demo envoy in a container and share volume to /tmp
### Host to container communication with UDS
```curl --unix-socket /tmp/demo_listener.sock any/demo/server/ok``` What's the difference between this command and the command in Experiment 3? 
Right now we are query on host UDS, and this host UDS is binded to demo envoy container UDS.
### Container to container communication with UDS
```curl --unix-socket /tmp/demo_listener.sock any/appnet/agent/status``` You will see connection refused error. Demo envoy container is trying 
to access /tmp/appnet_admin.sock which is created by appnet-agent container. We need to set the 'other permission'. Interestingly, demo envoy container 
needs other write permission.

```sudo chmod 756 /tmp/appnet_admin.sock``` and run above command again. You will get status information about Appnet envoy proxy.
