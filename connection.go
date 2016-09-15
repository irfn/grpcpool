package grpcpool

import (
	"google.golang.org/grpc"
)

type Connection interface {
	Close() error
}

type GrpcConnection struct {
	pool     Pool
	GrpcConn *grpc.ClientConn
}

func (self *GrpcConnection) Close() error {
	return self.pool.put(self.GrpcConn)
}
