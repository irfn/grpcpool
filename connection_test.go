package grpcpool

import (
	"testing"

	"google.golang.org/grpc"

	"github.com/stretchr/testify/mock"
)

type MockPool struct {
	mock.Mock
	mockConnection Connection
}

type MockConnection struct {
	mock.Mock
}

func (m MockPool) Get() (Connection, error) {
	return nil, nil
}

func (m MockPool) Close() {

}

func (m MockPool) put(client *grpc.ClientConn) error {
	m.Called(client)
	return nil
}

func TestShouldPutConnectionBackInPoolAfterClose(t *testing.T) {
	mockPool := new(MockPool)

	connection := &GrpcConnection{pool: mockPool, GrpcConn: &grpc.ClientConn{}}
	mockPool.On("put", connection.GrpcConn).Return(nil)
	connection.Close()
}
