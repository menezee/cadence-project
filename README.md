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

❯ make bins
❯ make worker

# in a different terminal
❯ alias cadence="docker run --rm ubercadence/cli:master --address host.docker.internal:7933 "
❯ make start

```

### TODO List
- [ ] HTTP Signal
- [ ] Makefile
- [ ] New Activity
