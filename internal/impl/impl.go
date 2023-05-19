package impl

import (
	"context"
	"runtime/debug"
	demoProto "template/demo"
	"time"

	log "github.com/sirupsen/logrus"
)

type DemoService struct {
	demoProto.UnimplementedDemoServiceServer
}

func NewDemoService() *DemoService {
	return &DemoService{}
}

func (s *DemoService) OneWay(ctx context.Context, req *demoProto.ReqPkg) (res *demoProto.RespPkg, err error) {
	CurrentReqCount.Inc()
	defer func() {
		CurrentReqCount.Dec()
		if err := recover(); err != nil {
			log.Printf("Work panic with %s %s\n", err, string(debug.Stack()))
		}
	}()
	res, err = &demoProto.RespPkg{Code: req.Age, Msg: "hello " + req.Name}, nil
	return
}

func (s *DemoService) HalfStream(req *demoProto.ReqPkg, stream demoProto.DemoService_HalfStreamServer) error {
	CurrentReqCount.Inc()
	defer func() {
		CurrentReqCount.Dec()
		if err := recover(); err != nil {
			log.Printf("Work panic with %s %s\n", err, string(debug.Stack()))
		}
	}()
	log.Println("half stream rev: ", req.Age, req.Name)
	for {
		if err := stream.Send(&demoProto.RespPkg{Code: 303, Msg: "continue"}); err != nil {
			log.Println("stream send error:", err)
		}
		if err := stream.Send(&demoProto.RespPkg{Code: 200, Msg: "success"}); err != nil {
			log.Println("stream send error:", err)
		}
		break
	}
	return nil
}

func (s *DemoService) Stream(stream demoProto.DemoService_StreamServer) error {
	CurrentReqCount.Inc()
	defer func() {
		CurrentReqCount.Dec()
		if err := recover(); err != nil {
			log.Printf("Work panic with %s %s\n", err, string(debug.Stack()))
		}
	}()

	for {
		req, err := stream.Recv()
		log.Println("stream rev: ", req.Age, req.Name)
		time.Sleep(time.Millisecond * 300)
		if err != nil {
			log.Println("stream rev error:", err)
			break
		}
		if req.Age == 18 {
			if err := stream.Send(&demoProto.RespPkg{Code: 200, Msg: "success"}); err != nil {
				log.Println("stream send error:", err)
			}
			break
		} else {
			if err := stream.Send(&demoProto.RespPkg{Code: 303, Msg: "continue"}); err != nil {
				log.Println("stream send error:", err)
			}
		}
	}
	return nil
}
