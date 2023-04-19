# go-service-template

大家都会维护接口文档? 开发接口服务吧? 你会碰到已下的场景么?

* 你希望接口文档在git仓库里面管理么, 这样你可以看到接口文档的历史版本
* 你希望通过接口文档直接生成接口代码么, 包括序列化对象那些
* 你希望你的接口文档可以跟着服务走么, 直接通过服务暴露的端口在浏览器查看么
* 你希望定义的接口可以同时被内部通过rpc调用, 也可以被外部通过http+json调用么
* 你希望定义的流式接口可以同时被内部通过rpc调用, 也可以被外部通过http-trunk+json调用么
* 你希望定义的接口可以直接被websocket调用么
* 你希望你的接口可以暴露出一些默认指标给Prometheus监控么
* 你希望你的服务可以自带Profile接口, 可以通过Go Profile工具来分析你的服务么
* 你希望你的服务自带Gops, 可以通过Gops工具来分析你的服务么
* 你希望你的服务可以自带退出信号的捕获和优雅退出么
* 你希望你的服务可以自带日志的配置文件么, 可以配置日志级别、输出的文件、格式、文件回
* 你希望你的服务支持交叉编译, 生成不同平台和架构的二进制文件么
* 你希望使用配置文件管理程序的一些行为么
* 你需要多阶段构建, 把你的服务打包成一个合理体积的 docker image么

如果你需要上面的任何一个功能, 你可以了解下这个项目。

## 依赖说明

使用了 `grep-gateway` `protoc-gen-doc` `grpcurl` `gops` `protoc` 工具不了解的可以看下如何使用

## 快速开始

1. 在 `demo` 目录中定义你的接口, 例如 `demo.proto`
2. 使用 `demo/gen.sh` 生成接口代码
3. 在 `internal/impl` 目录中实现你的接口, 例如 `impl.go`
4. 在 `conf/services.yaml` 中配置服务暴露的端口以及日志信息

## 使用说明

### 接口定义

本项目采用grpc protobuf作为接口定义语言，使用grpc-gateway作为grpc和http的转换桥梁。

所以你需要熟悉grpc接口定义语言, 如果你不熟悉, 你可以化个30分钟的时间看看官方的文档

[grpc接口定义语言](https://developers.google.com/protocol-buffers/docs/proto3)

接口定义语言定义好之后, 你需要使用protoc工具生成api代码和api文档(html 当然也有其他的格式)

如果你想生成其他格式的接口文档, 请参阅 [protoc-gen-doc](https://github.com/pseudomuto/protoc-gen-doc)

可以使用浏览器查看接口文档。

http://localhost:8070/api/docs

### 配置管理

本项目使用yaml作为配置文件的格式, 你可以在 `conf/services.yaml` 中配置服务的端口, 日志的配置文件路径等

当然你可以扩展其他的配置信息, 比如业务逻辑上的一些开关, 你可以修改 `internal/config/config.go` 中的 `Config` 结构体

### 日志管理

本项目使用logrus作为日志框架, 你可以在 `conf/services.yaml` 中配置日志的输出级别, 输出文件, 输出格式, 文件回滚配置等

如果你想配置额外的信息，比如把日志输出到syslogd等，可以参考logrus的文档

### 优雅退出

服务启动之后，会监听系统的退出信号，当接收到退出信号之后，会等当前的请求处理完毕之后，再退出服务。

当然这需要你的业务逻辑支持超时退出，你需要在处理函数中加入
    
```go
CurrentReqCount.Inc()
defer func() {
    CurrentReqCount.Dec()
}()
```

### 指标监控

自动集成了prometheus暴露基础指标, 可以使用下面的地址查看

http://localhost:8060/metrics

### Profile

支持gops来debug服务

[gops使用](https://github.com/google/gops)

### 交叉编译

提供了默认的 `build.sh` 支持不同系操作系统和不同架构的编译

## 代码自测

这里提供了一些通过命令行调用服务的工具, 只用这些工具就可以测试你的服务.

同时你也很容易的使用一些性能测试工具, 比如apache ab, ghz等

当然你也可以直接使用ide调用你的服务, 这里提供了一个 `Goland` 的示例。

所有的测试工具都在 `testing` 目录下有相关的参考。

### curl

```shell
curl -X POST -d '{"age":18,"name":"scguo"}' http://127.0.0.1:8090/demo
curl -X POST -d '{"name":"scguo","age":1}{"name":"scguo","age":18}' -H "Transfer-Encoding: chunked" 127.0.0.1:8090/stream
```

### grpCurl
```shell
grpcurl -plaintext -d '{"age":18,"name":"scguo"}' 127.0.0.1:8080 example.DemoService/OneWay
grpcurl -plaintext -d @ 127.0.0.1:8080  example.DemoService.Stream < grpc_stream.dat
```

### websocket

```shell
websocat "ws://127.0.0.1:8091/stream?method=POST" -H='Origin: ws://127.0.0.1:8091/stream?method=POST' -v
```

## todo

1. 使用日志context把一次请求的所有日志都打印到一起
