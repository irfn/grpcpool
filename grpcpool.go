package grpcpool

import (
	"sync"

	"google.golang.org/grpc"
)

func newConnectionPool(activeCount int, dialFunc func() (*grpc.ClientConn, error)) (*ConnectionPool, error) {
	pool := &ConnectionPool{mu: sync.Mutex{}}
	for i := 0; i < activeCount; i++ {
		client, error := dialFunc()
		if error != nil {
			pool.Close()
			return nil, error
		}
		pool.put(client)
	}
	return pool, nil
}
