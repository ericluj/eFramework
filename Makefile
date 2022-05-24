run:
	go mod tidy
	go mod vendor
	go run main.go

# 编译proto生成pb.go文件
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