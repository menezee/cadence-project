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
❯ go build -o bins/worker main.go simple-activity.go simple-workflow.go
❯ ./bins/worker

# in a different terminal
❯ cadence --domain tdc workflow run --tl tdcTasks --wt main.TDCWorkflow --et 60 -i '"cadence"'
```

### TODO List
- [ ] HTTP Signal
- [ ] Makefile
- [ ] New Activity
