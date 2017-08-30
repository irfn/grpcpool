package grpcpool

import (
	"testing"

	"google.golang.org/grpc"

	"github.com/stretchr/testify/mock"
	"context"
	"time"
	"net"
)

type MockPool struct {
	mock.Mock
	mockConnection Connection
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

	mockConnection := &MockConnection{}
	conn, _ := grpc.DialContext(context.Background(), "fakeaddr", grpc.WithDialer(func(string, time.Duration) (net.Conn, error) {
		return mockConnection, nil
	}), grpc.WithInsecure())

	grpcConnection := &GrpcConnection{pool: mockPool, GrpcConn: conn}
	mockPool.On("put", grpcConnection.GrpcConn).Return(nil)

	grpcConnection.Close()

	mockConnection.AssertNotCalled(t, "Close")
}

func TestShouldPutEvictConnectionFromPool(t *testing.T) {
	mockPool := new(MockPool)

	mockConnection := &MockConnection{}
	conn, _ := grpc.DialContext(context.Background(), "fakeaddr", grpc.WithDialer(func(string, time.Duration) (net.Conn, error) {
		return mockConnection, nil
	}), grpc.WithInsecure())

	grpcConnection := &GrpcConnection{pool: mockPool, GrpcConn: conn}
	mockPool.On("put", grpcConnection.GrpcConn).Return(nil)
	mockConnection.On("Close").Return(nil)

	grpcConnection.Evict()

	mockConnection.AssertExpectations(t)
}
