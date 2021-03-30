package proxy

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// TcpServer tcp server common
type TcpServer struct {
	Server
	host string
	lock sync.Mutex

	listener net.Listener

	tcpConn atomic.Value
	handle  func(net.Conn, func(string) (net.Conn, error))
}

type Option struct {
	TcpConn func(string) (net.Conn, error)
}

// NewTCPServer create new TCP listener
func NewTCPServer(host string, handle func(net.Conn, func(string) (net.Conn, error)), modeOption ...func(*Option)) (TCPServer, error) {
	if host == "" {
		return nil, errors.New("host empty")
	}
	if handle == nil {
		return nil, errors.New("handle is must")
	}
	o := &Option{
		TcpConn: func(s string) (net.Conn, error) {
			return net.DialTimeout("tcp", s, 20*time.Second)
		},
	}
	for index := range modeOption {
		if modeOption[index] == nil {
			continue
		}
		modeOption[index](o)
	}

	s := &TcpServer{
		host:    host,
		handle:  handle,
		tcpConn: atomic.Value{},
	}
	s.tcpConn.Store(o.TcpConn)

	err := s.run()
	if err != nil {
		return nil, fmt.Errorf("server Run -> %v", err)
	}
	return s, nil
}

func (t *TcpServer) UpdateListen(host string) (err error) {
	if t.host == host {
		return
	}
	_ = t.Close()

	t.lock.Lock()
	defer t.lock.Unlock()

	if host == "" {
		return
	}

	t.host = host

	fmt.Println("UpdateListen create new server")
	return t.run()
}

func (t *TcpServer) SetTCPConn(conn func(string) (net.Conn, error)) {
	if conn == nil {
		return
	}
	t.tcpConn.Store(conn)
}

func (t *TcpServer) getTCPConn() func(string) (net.Conn, error) {
	y, ok := t.tcpConn.Load().(func(string) (net.Conn, error))
	if ok {
		return y
	}
	return func(s string) (net.Conn, error) {
		return net.Dial("tcp", s)
	}
}

func (t *TcpServer) GetListenHost() string {
	return t.host
}

func (t *TcpServer) run() (err error) {
	fmt.Println("New TCP Server:", t.host)
	t.listener, err = net.Listen("tcp", t.host)
	if err != nil {
		return fmt.Errorf("TcpServer:run() -> %v", err)
	}

	go t.process()
	return
}

func (t *TcpServer) process() error {
	t.lock.Lock()
	defer t.lock.Unlock()

	var tempDelay time.Duration
	for {
		c, err := t.listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				log.Printf("checked tcp server closed: %v\n", err)
				return fmt.Errorf("checked tcp server closed: %v", err)
			}

			// from https://golang.org/src/net/http/server.go?s=93655:93701#L2977
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}

				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}

				log.Printf("tcp sever: Accept error: %v; retrying in %v\n", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}

			log.Printf("tcp server accept failed: %v\n", err)
			return fmt.Errorf("tcp server accept failed: %v", err)
		}

		tempDelay = 0

		go func() {
			defer c.Close()
			t.handle(c, t.getTCPConn())
		}()
	}
}

func (t *TcpServer) Close() error {
	if t.listener == nil {
		return nil
	}
	return t.listener.Close()
}
