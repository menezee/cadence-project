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
❯ curl -X POST "http://localhost:3030/bank-message-signal?workflowId=<workflowIdFromPreviousCommand>&bankMessage=tudo%20certo"
```
 
To start Order Workflow:
```
# skip the docker run and make start
❯ make http-server
❯ curl -X POST "http://localhost:3030/create-order?totalValue=50"
❯ curl -X GET "http://localhost:3030/get-status?workflowId=<workflowIdFromPreviousPOSTResponse>"
❯ curl -X POST "http://localhost:3030/order-received?workflowId=<workflowIdFromPreviousPOSTResponse>"
``` 

