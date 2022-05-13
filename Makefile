run:
	go run main.go

# 编译proto生成pb.go文件
gen:
	rm -rf rpc/*
	protoc --go_out=./rpc --go-grpc_out=./rpc ./proto/*/*.proto