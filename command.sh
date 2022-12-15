# 生成协议
protoc -I=demo --doc_out=demo --doc_opt=html,index.html --go_out=apiproto --go_opt=paths=source_relative --go-grpc_out=demo --go-grpc_opt=paths=source_relative --grpc-gateway_out=demo --grpc-gateway_opt=paths=source_relative ./demo/demo.proto

# build
docker build -t scguo/mydemo:latest .

# run
docker run -p 8080:8080 -p 8090:8090 -p 8070:8070 -d scguo/mydemo:latest