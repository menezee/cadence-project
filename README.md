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

### TODO List
- [ ] HTTP Signal
- [ ] Makefile
- [ ] New Activity
