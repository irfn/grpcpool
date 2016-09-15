package grpcpool

import (
	"container/list"
	"errors"
	"sync"

	"google.golang.org/grpc"
)

type Pool interface {
	Get() (Connection, error)
	Close()
	put(client *grpc.ClientConn) error
}

type ConnectionPool struct {
	mu     sync.Mutex
	cond   *sync.Cond
	closed bool
	idle   list.List
}

func (self *ConnectionPool) Get() (Connection, error) {
	conn, err := self.get()
	if err != nil {
		// TODO : error handling
		return nil, err
	}
	return &GrpcConnection{pool: self, GrpcConn: conn}, nil
}

func (self *ConnectionPool) Close() {
	self.mu.Lock()
	idle := self.idle
	self.idle.Init()
	self.closed = true
	if self.cond != nil {
		self.cond.Broadcast()
	}
	self.mu.Unlock()
	for e := idle.Front(); e != nil; e = e.Next() {
		e.Value.(*grpc.ClientConn).Close()
	}
}

func (self *ConnectionPool) put(client *grpc.ClientConn) error {
	self.mu.Lock()
	if self.closed {
		self.mu.Unlock()
		return client.Close()
	}

	self.idle.PushFront(client)
	self.mu.Unlock()
	if self.cond != nil {
		self.cond.Signal()
	}
	return nil
}

func (self *ConnectionPool) get() (*grpc.ClientConn, error) {
	self.mu.Lock()
	for {
		element := self.idle.Front()

		if element != nil {
			self.idle.Remove(element)
			client := element.Value.(*grpc.ClientConn)
			self.mu.Unlock()
			return client, nil
		}

		if self.closed {
			return nil, errors.New("Pool is closed")
		}

		// wait for client to be available
		if self.cond == nil {
			self.cond = sync.NewCond(&self.mu)
		}
		self.cond.Wait()
	}
}
