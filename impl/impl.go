package impl

import (
	"context"
	demoProto "template/demo"

	log "github.com/sirupsen/logrus"
)

type DemoService struct {
	demoProto.UnimplementedDemoServiceServer
}

func NewDemoService() *DemoService {
	return &DemoService{}
}

func (s *DemoService) OneWay(ctx context.Context, req *demoProto.ReqPkg) (*demoProto.RespPkg, error) {
	log.WithField("req", req).Info("OneWay")
	return &demoProto.RespPkg{Code: req.Age, Msg: "hello " + req.Name}, nil
}
