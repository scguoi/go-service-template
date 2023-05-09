package main

import (
	"context"
	"fmt"
	"github.com/google/gops/agent"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
	"net/http/pprof"
	"template/internal/config"
	"template/internal/gracefulstop"
	"template/internal/impl"

	demoProto "template/demo"

	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/tmc/grpc-websocket-proxy/wsproxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	_ "net/http/pprof"
	_ "template/internal/config"
	_ "template/internal/logc"
)

func main() {
	// grpc server
	log.Println("hi service is starting...")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.ServiceConfig.GRPCPort))
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpcPrometheus.StreamServerInterceptor),
		grpc.UnaryInterceptor(grpcPrometheus.UnaryServerInterceptor),
	)
	demoProto.RegisterDemoServiceServer(grpcServer, impl.NewDemoService())
	grpcPrometheus.Register(grpcServer)
	reflection.Register(grpcServer)
	log.Println("Serving gRPC on", fmt.Sprintf(":%d", config.ServiceConfig.GRPCPort))
	go func() { log.Fatalln(grpcServer.Serve(lis)) }()

	// grpc gateway
	conn, err := grpc.DialContext(
		context.Background(),
		fmt.Sprintf("0.0.0.0:%d", config.ServiceConfig.GRPCPort),
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

	// grpc websocket
	gwWSServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.ServiceConfig.WSPort),
		Handler: wsproxy.WebsocketProxy(gwMux),
	}
	log.Println("Serving gRPC-Websocket on", fmt.Sprintf(":%d", config.ServiceConfig.WSPort))
	go func() { log.Fatalln(gwWSServer.ListenAndServe()) }()

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

	// prometheus
	go func() {
		log.Println("Serving Prometheus on", fmt.Sprintf(":%d", config.ServiceConfig.MetricPort))
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
		_ = http.ListenAndServe(fmt.Sprintf(":%d", config.ServiceConfig.MetricPort), mux)
	}()

	// gOps
	agentOptions := agent.Options{
		ShutdownCleanup: true,
		Addr:            fmt.Sprintf(":%d", config.ServiceConfig.AgentOpsPort),
	}
	if err := agent.Listen(agentOptions); err != nil {
		log.Fatalln(err)
	}

	// graceful stop
	signalChan := gracefulstop.NewShutdownSignal()
	gracefulstop.WaitExit(signalChan, func(ctx context.Context) {
		err := gwServer.Shutdown(ctx)
		if err != nil {
			log.Println("gwServer shutdown failed", err)
		} else {
			log.Println("gwServer shutdown succeed")
		}
		err = gwWSServer.Shutdown(ctx)
		if err != nil {
			log.Println("gwWSServer shutdown failed", err)
		} else {
			log.Println("gwWSServer shutdown succeed")
		}
		grpcServer.GracefulStop()
		log.Println("grpc server graceful stop")
		err = apiSrv.Shutdown(ctx)
		if err != nil {
			log.Println("apiSrv shutdown failed", err)
		} else {
			log.Println("apiSrv shutdown succeed")
		}
		agent.Close()
	})
}
