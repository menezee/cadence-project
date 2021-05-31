.PHONY : bins

bins:
	go build -o bins/worker *.go

worker:
	./bins/worker

start:
	docker run --rm ubercadence/cli:master --address host.docker.internal:7933 --domain tdc workflow run --tl tdcTasks --wt main.TDCWorkflow --et 60 -i '"cadence"'

signal:
	docker run --rm ubercadence/cli:master --address host.docker.internal:7933 --domain tdc workflow signal -w $(wf) -n signalForTDC -i '25'

register-domain:
	docker run --rm ubercadence/cli:master --address host.docker.internal:7933 --domain tdc domain register
