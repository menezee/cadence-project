# Cadence TDC

### Requirements
- [Docker]()
- [Go]()

### Instructions
```
❯ cd docker
❯ docker-compose up

# in a different terminal
❯ go build -i -o bins/worker main.go simple-activity.go simple-workflow.go
❯ ./bins/worker

# in a different terminal
❯ cadence --domain tdc workflow run --tl tdcTasks --wt main.TDCWorkflow --et 60 -i '"cadence"'
```