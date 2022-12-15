grpcurl -plaintext -d '{"age":18,"name":"scguo"}' 127.0.0.1:8080 example.DemoService/OneWay

curl -X POST -d '{"age":18,"name":"scguo"}' http://127.0.0.1:8090/demo