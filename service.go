package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	_ "net/http/pprof"

	"template/internal/config"
	"template/internal/gracefulstop"
	"template/internal/impl"
	"template/internal/logc"

	"github.com/google/gops/agent"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	demoProto "template/protocol"

	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/scguoi/grpc-gateway/v2/runtime"
	log "github.com/sirupsen/logrus"
	"github.com/tmc/grpc-websocket-proxy/wsproxy"
	"google.golang.org/grpc"
)

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	serverCert, err := tls.LoadX509KeyPair("conf/cert.pem", "conf/key.pem")
	if err != nil {
		return nil, err
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
	}

	return credentials.NewTLS(tlsConfig), nil
}

func main() {
	// load config and initial
	config.Load()
	logc.Initial()
	impl.Initial()

	// grpc server
	log.Printf("hi service is starting...")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.ServiceConfig.GRPCPort))
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}
	var grpcServer *grpc.Server
	var cert credentials.TransportCredentials
	var conn *grpc.ClientConn

	if config.ServiceConfig.UsingGrpcs {
		cert, err = loadTLSCredentials()
		if err != nil {
			log.Fatalln("load tls cert failed:", err)
		}
		grpcServer = grpc.NewServer(
			grpc.Creds(cert),
			grpc.StreamInterceptor(grpcPrometheus.StreamServerInterceptor),
			grpc.UnaryInterceptor(grpcPrometheus.UnaryServerInterceptor),
		)
	} else {
		grpcServer = grpc.NewServer(
			grpc.StreamInterceptor(grpcPrometheus.StreamServerInterceptor),
			grpc.UnaryInterceptor(grpcPrometheus.UnaryServerInterceptor),
		)
	}

	demoProto.RegisterDemoServiceServer(grpcServer, impl.NewDemoService())
	grpcPrometheus.Register(grpcServer)
	reflection.Register(grpcServer)
	log.Printf("Serving gRPC on %s", fmt.Sprintf(":%d", config.ServiceConfig.GRPCPort))
	go func() { log.Fatalln(grpcServer.Serve(lis)) }()

	// grpc gateway
	if config.ServiceConfig.UsingGrpcs {
		conn, err = grpc.DialContext(
			context.Background(),
			fmt.Sprintf("127.0.0.1:%d", config.ServiceConfig.GRPCPort),
			grpc.WithBlock(),
			grpc.WithTransportCredentials(cert),
		)

		if err != nil {
			log.Fatalln("Failed to dial server:", err)
		}
	} else {
		conn, err = grpc.DialContext(
			context.Background(),
			fmt.Sprintf("127.0.0.1:%d", config.ServiceConfig.GRPCPort),
			grpc.WithBlock(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			log.Fatalln("Failed to dial server:", err)
		}
	}

	gwMux := runtime.NewServeMux()
	err = demoProto.RegisterDemoServiceHandler(context.Background(), gwMux, conn)
	gwServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.ServiceConfig.HTTPPort),
		Handler: gwMux,
	}
	log.Printf("Serving gRPC-Gateway on %s", fmt.Sprintf(":%d", config.ServiceConfig.HTTPPort))
	go func() { log.Fatalln(gwServer.ListenAndServe()) }()

	// grpc websocket
	gwWSServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.ServiceConfig.WSPort),
		Handler: wsproxy.WebsocketProxy(gwMux),
	}
	log.Printf("Serving gRPC-Websocket on %s", fmt.Sprintf(":%d", config.ServiceConfig.WSPort))
	go func() { log.Fatalln(gwWSServer.ListenAndServe()) }()

	// api docs
	router := httprouter.New()
	router.GET("/api/docs", demoProto.APIProto)
	bind := fmt.Sprintf(":%d", config.ServiceConfig.APIPort)
	log.Printf("Serving API Proto starting on %s", bind)
	apiSrv := &http.Server{
		Addr:    bind,
		Handler: router,
	}
	go func() { log.Fatalln(apiSrv.ListenAndServe()) }()

	// prometheus
	go func() {
		log.Printf("Serving Prometheus on %s", fmt.Sprintf(":%d", config.ServiceConfig.MetricPort))
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
			log.Printf("gwServer shutdown failed %v", err)
		} else {
			log.Printf("gwServer shutdown succeed")
		}
		err = gwWSServer.Shutdown(ctx)
		if err != nil {
			log.Printf("gwWSServer shutdown failed %v", err)
		} else {
			log.Printf("gwWSServer shutdown succeed")
		}
		grpcServer.GracefulStop()
		log.Printf("grpc server graceful stop")
		err = apiSrv.Shutdown(ctx)
		if err != nil {
			log.Printf("apiSrv shutdown failed %v", err)
		} else {
			log.Printf("apiSrv shutdown succeed")
		}
		agent.Close()
	})
}
