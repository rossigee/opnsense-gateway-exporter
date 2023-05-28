image:
	docker build . -t rossigee/opnsense-gateway-exporter:test

push: image
	docker push rossigee/opnsense-gateway-exporter:test

all:
	go build -o main
