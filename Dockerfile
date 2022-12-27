FROM golang:1.19.4-bullseye AS builder
WORKDIR /DEMO
COPY . .
ENV GOPROXY=https://goproxy.cn,direct
RUN go build -o demo_service service.go

FROM busybox as prod
COPY --from=builder /DEMO/demo_service /DEMO/demo_service
COPY --from=builder /DEMO/conf /DEMO/conf
COPY --from=builder /DEMO/apiproto /DEMO/apiproto
EXPOSE 8080
EXPOSE 8090
EXPOSE 8070
EXPOSE 8060
WORKDIR /DEMO
CMD ["./demo_service"]