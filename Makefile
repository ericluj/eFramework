test:
	echo "test"

# 编译proto生成pb.go文件
gen:
	rm -rf rpc/*
	protoc --go-grpc_out=./rpc ./proto/*/*.proto