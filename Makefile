.PHONY : bins http-server

bins:
	go build -o bins/worker *.go & go build -o bins/http-server http-server/main.go

worker:
	./bins/worker

http-server:
	./bins/http-server

register-domain:
	docker run --rm ubercadence/cli:master --address host.docker.internal:7933 --domain tdc domain register

start:
	docker run --rm ubercadence/cli:master --address host.docker.internal:7933 --domain tdc workflow run --tl tdcTasks --wt github.com/menezee/cadence-project/eats.CourierTipWorkflow --et 60 -i '50'

signal:
	docker run --rm ubercadence/cli:master --address host.docker.internal:7933 --domain tdc workflow signal -w $(wf) -n BankConfirmationSignalToken -i '"Everything done from the bank perspective!"'
