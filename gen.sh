protoc -I=demo --go_out=demo --go_opt=paths=source_relative --go-grpc_out=demo --go-grpc_opt=paths=source_relative --grpc-gateway_out=demo --grpc-gateway_opt=paths=source_relative ./demo/demo.proto
