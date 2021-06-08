# Cadence TDC

### Requirements
- [Docker]()
- [Go]()

### Instructions
```
❯ cd docker
❯ docker-compose up

# in a different terminal
❯ go get go.uber.org/cadence
# or go mod download

❯ make register-domain
❯ make bins
❯ make worker

# in a different terminal
❯ alias cadence="docker run --rm ubercadence/cli:master --address host.docker.internal:7933 "
❯ make start

# in a different terminal
❯ make signal wf=<workflowIdFromPreviousCommand>
```

if you want to send the signal through an http request:
```shell
❯ make http-server
❯ curl -X POST "http://localhost:3030/signal-workflow?workflowId=<workflowIdFromPreviousCommand>&age=25"
```

### TODO List
- [ ] HTTP Signal
- [ ] Makefile
- [ ] New Activity
