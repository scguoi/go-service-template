# 生成协议
protoc -I=demo --go_out=demo --go_opt=paths=source_relative --go-grpc_out=demo --go-grpc_opt=paths=source_relative --grpc-gateway_out=demo --grpc-gateway_opt=paths=source_relative ./demo/demo.proto

# build
docker build -t scguo/mydemo:latest .

# run
docker run -p 8080:8080 -p 8090:8090 -d scguo/mydemo:latest