package grpcpool

import (
	"github.com/stretchr/testify/mock"
	"time"
	"net"
	"golang.org/x/net/http2"
)

type MockConnection struct {
	mock.Mock
}

func (m *MockConnection) Read(b []byte) (n int, err error) {
	return 0, nil
}

func (m *MockConnection) Write(b []byte) (n int, err error) {
	return len(http2.ClientPreface), nil
}

func (m *MockConnection) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockConnection) LocalAddr() net.Addr {
	args := m.Called()
	return args[0].(net.Addr)
}

func (m *MockConnection) RemoteAddr() net.Addr {
	args := m.Called()
	return args[0].(net.Addr)
}

func (m *MockConnection) SetDeadline(t time.Time) error {
	args := m.Called(t)
	return args.Error(0)
}

func (m *MockConnection) SetReadDeadline(t time.Time) error {
	args := m.Called(t)
	return args.Error(0)
}

func (m *MockConnection) SetWriteDeadline(t time.Time) error {
	args := m.Called(t)
	return args.Error(0)
}
