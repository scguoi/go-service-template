grpcurl -plaintext -d '{"age":18,"name":"scguo"}' 127.0.0.1:8080 example.DemoService/OneWay
grpcurl -plaintext -d @ 127.0.0.1:8080  example.DemoService.Stream < grpc_stream.dat