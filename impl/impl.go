package impl

import (
	"context"
	"runtime/debug"
	demoProto "template/demo"

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
	log.WithField("req", req).Info("OneWay")
	panic("panic testing")
}
