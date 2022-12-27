# go-service-template
go 服务模板 支持grpc和http请求，使用logrus配置日志，使用grpc和grpc-gateway暴露服务。

使用yaml配置文件配置服务。

可以使用浏览器查看接口文档。

http://localhost:8070/api/docs

支持退出信号的捕获和优雅退出。

## gen.sh
用于生成grpc grpc-gateway的代码，本地需要安装protoc protoc-gen-go protoc-gen-grpc-gateway

## logc
logrus的配置文件，可以配置日志的输出级别，输出文件，输出格式，文件回滚配置等

## config
配置文件，可以配置服务的端口，grpc的端口，grpc-gateway的端口，日志的配置文件路径等

## build.sh
使用go build编译生成二进制文件，支持交叉编译，输出不同的cpu架构和操作系统的二进制文件

## testing

```shell
grpcurl -plaintext -d '{"age":18,"name":"scguo"}' 127.0.0.1:8080 example.DemoService/OneWay

curl -X POST -d '{"age":18,"name":"scguo"}' http://127.0.0.1:8090/demo
```

## todo

1. 使用日志context把一次请求的所有日志都打印到一起
2. 使用prometheus监控服务