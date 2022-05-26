mod:
	go mod tidy
	go mod vendor

gen:
	rm -rf rpc/*
	protoc -I . \
    	--go_out ./rpc \
    	--go-grpc_out ./rpc \
    	./proto/*/*.proto
	protoc -I . \
		--grpc-gateway_out ./rpc \
    	--grpc-gateway_opt logtostderr=true \
    	./proto/*/*.proto

image:
	docker build -t front:v1 -f deploy/front/Dockerfile .
	docker build -t sample:v1 -f deploy/sample/Dockerfile .

up:
	docker-compose up -d

stop:
	docker-compose stop

recreate:
	docker-compose up -d --force-recreate

sample:
	go run cmd/sample/main.go

front:
	go run cmd/front/main.go