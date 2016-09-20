package grpcpool

import (
	"google.golang.org/grpc"
)

type Connection interface {
	Close() error
	Get() *grpc.ClientConn
}

type GrpcConnection struct {
	pool     Pool
	GrpcConn *grpc.ClientConn
}

func (self *GrpcConnection) Close() error {
	return self.pool.put(self.GrpcConn)
}

func (self *GrpcConnection) Get() *grpc.ClientConn {
	return self.GrpcConn
}
