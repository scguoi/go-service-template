package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"template/config"
	"template/impl"
	"template/osutils"

	demoProto "template/demo"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	_ "template/config"
	_ "template/logc"
)

func main() {
	// grpc server
	log.Debugln("hi service is starting...")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.ServiceConfig.GRPCPort))
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	demoProto.RegisterDemoServiceServer(grpcServer, impl.NewDemoService())
	log.Println("Serving gRPC on", fmt.Sprintf(":%d", config.ServiceConfig.GRPCPort))
	go func() { log.Fatalln(grpcServer.Serve(lis)) }()

	// grpc gateway
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

	// api docs
	router := httprouter.New()
	router.GET("/api/docs", demoProto.APIProto)
	bind := fmt.Sprintf(":%d", config.ServiceConfig.APIPort)
	log.Println("Serving API Proto starting on", bind)
	apiSrv := &http.Server{
		Addr:    bind,
		Handler: router,
	}
	go func() { log.Fatalln(apiSrv.ListenAndServe()) }()

	// graceful stop
	signalChan := osutils.NewShutdownSignal()
	osutils.WaitExit(signalChan, func(ctx context.Context) {
		err := gwServer.Shutdown(ctx)
		if err != nil {
			log.Println("gwServer shutdown failed", err)
		} else {
			log.Println("gwServer shutdown succeed")
		}
		grpcServer.GracefulStop()
		log.Println("grpc server graceful stop")
		err = apiSrv.Shutdown(ctx)
		if err != nil {
			log.Println("apiSrv shutdown failed", err)
		} else {
			log.Println("apiSrv shutdown succeed")
		}
	})
}
