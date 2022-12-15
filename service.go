package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"template/config"
	"template/impl"

	demoProto "template/demo"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	_ "template/config"
	_ "template/logc"
)

func main() {
	log.Debugln("hi service is starting...")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.ServiceConfig.GRPCPort))
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}
	s := grpc.NewServer()
	demoProto.RegisterDemoServiceServer(s, impl.NewDemoService())
	log.Println("Serving gRPC on", fmt.Sprintf(":%d", config.ServiceConfig.GRPCPort))
	go func() { log.Fatalln(s.Serve(lis)) }()

	conn, err := grpc.DialContext(
		context.Background(),
		"0.0.0.0:8080",
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}
	gwMux := runtime.NewServeMux()
	err = demoProto.RegisterDemoServiceHandler(context.Background(), gwMux, conn)
	gwServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.ServiceConfig.HTTPPort),
		Handler: gwMux,
	}
	log.Println("Serving gRPC-Gateway on", fmt.Sprintf(":%d", config.ServiceConfig.HTTPPort))
	go func() { log.Fatalln(gwServer.ListenAndServe()) }()

	router := httprouter.New()
	router.GET("/api/docs", demoProto.APIProto)
	bind := fmt.Sprintf(":%d", config.ServiceConfig.APIPort)
	log.Println("Serving API Proto starting on", bind)
	log.Fatalln(http.ListenAndServe(bind, router))
}
