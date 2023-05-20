package impl

import (
	"context"
	"runtime/debug"
	demoProto "template/demo"
	"template/internal/logc"
	"time"
)

type DemoService struct {
	demoProto.UnimplementedDemoServiceServer
}

func NewDemoService() *DemoService {
	return &DemoService{}
}

type service struct {
	log *logc.BizLog
}

func (s *DemoService) OneWay(ctx context.Context, req *demoProto.ReqPkg) (res *demoProto.RespPkg, err error) {
	return (&service{log: logc.NewBizLog()}).OneWay(ctx, req)
}

func (s *DemoService) HalfStream(req *demoProto.ReqPkg, stream demoProto.DemoService_HalfStreamServer) error {
	return (&service{log: logc.NewBizLog()}).HalfStream(req, stream)
}

func (s *DemoService) Stream(stream demoProto.DemoService_StreamServer) error {
	return (&service{log: logc.NewBizLog()}).Stream(stream)
}

func (s *service) OneWay(ctx context.Context, req *demoProto.ReqPkg) (res *demoProto.RespPkg, err error) {
	CurrentReqCount.Inc()
	startTime := time.Now()
	defer func() {
		CurrentReqCount.Dec()
		ResponseTime.WithLabelValues("OneWay").Observe(time.Since(startTime).Seconds())
		if err := recover(); err != nil {
			s.log.Printf("Work panic with %s %s\n", err, string(debug.Stack()))
		}
		s.log.LoggerEnd()
	}()
	s.log.Print("oneway rev: ", req.Age, req.Name)
	res, err = &demoProto.RespPkg{Code: req.Age, Msg: "hello " + req.Name}, nil
	return
}

func (s *service) HalfStream(req *demoProto.ReqPkg, stream demoProto.DemoService_HalfStreamServer) error {
	CurrentReqCount.Inc()
	startTime := time.Now()
	defer func() {
		CurrentReqCount.Dec()
		ResponseTime.WithLabelValues("HalfStream").Observe(time.Since(startTime).Seconds())
		if err := recover(); err != nil {
			s.log.Printf("Work panic with %s %s\n", err, string(debug.Stack()))
		}
		s.log.LoggerEnd()
	}()
	s.log.Print("half stream rev: ", req.Age, req.Name)
	for {
		if err := stream.Send(&demoProto.RespPkg{Code: 303, Msg: "continue"}); err != nil {
			s.log.Printf("stream send error:", err)
		}
		time.Sleep(100 * time.Millisecond)
		if err := stream.Send(&demoProto.RespPkg{Code: 200, Msg: "success"}); err != nil {
			s.log.Printf("stream send error:", err)
		}
		break
	}
	return nil
}

func (s *service) Stream(stream demoProto.DemoService_StreamServer) error {
	CurrentReqCount.Inc()
	startTime := time.Now()
	defer func() {
		CurrentReqCount.Dec()
		ResponseTime.WithLabelValues("Stream").Observe(time.Since(startTime).Seconds())
		if err := recover(); err != nil {
			s.log.Printf("Work panic with %s %s\n", err, string(debug.Stack()))
		}
		s.log.LoggerEnd()
	}()

	for {
		req, err := stream.Recv()
		s.log.Print("stream rev: ", req.Age, req.Name)
		time.Sleep(time.Millisecond * 300)
		if err != nil {
			s.log.Printf("stream rev error:", err)
			break
		}
		if req.Age == 18 {
			if err := stream.Send(&demoProto.RespPkg{Code: 200, Msg: "success"}); err != nil {
				s.log.Printf("stream send error:", err)
			}
			break
		} else {
			if err := stream.Send(&demoProto.RespPkg{Code: 303, Msg: "continue"}); err != nil {
				s.log.Printf("stream send error:", err)
			}
		}
	}
	return nil
}
