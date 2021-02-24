package vmess

import (
	"crypto/tls"
	"io/ioutil"
	"net"
	"time"

	"github.com/lucas-clemente/quic-go"
)

func QuicDial(network, address string, port int, host string, certPath string) (net.Conn, error) {
	// conn, err := net.ListenUDP("udp")
	// if err != nil {
	// return nil, err
	// }

	var addr *net.UDPAddr
	var err error
	switch network {
	case "ip":
		addr = &net.UDPAddr{
			IP:   net.ParseIP(address),
			Port: port,
		}
	default:
		addr, err = net.ResolveUDPAddr("udp", address)
		if err != nil {
			return nil, err
		}
	}

	cert, err := ioutil.ReadFile(certPath)
	if err != nil {
		return nil, err
	}
	// key, err := ioutil.ReadFile(keyPath)
	// if err != nil {
	// 	return nil, err
	// }
	// certPair, err := tls.X509KeyPair(cert, key)
	// if err != nil {
	// 	return nil, err
	// }

	tlsConfig := &tls.Config{
		ServerName: host,
		Certificates: []tls.Certificate{
			{
				Certificate: [][]byte{cert},
			},
		},
	}

	quicConfig := &quic.Config{
		KeepAlive:          true,
		ConnectionIDLength: 12,
		HandshakeTimeout:   time.Second * 8,
		MaxIdleTimeout:     time.Second * 30,
	}

	conn, err := net.DialUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: 0}, addr)
	if err != nil {
		return nil, err
	}

	session, err := quic.Dial(conn, addr, host, tlsConfig, quicConfig)
	if err != nil {
		return nil, err
	}

	stream, err := session.OpenStream()
	if err != nil {
		return nil, err
	}

	return &interConn{
		stream: stream,
	}, nil
}

type interConn struct {
	net.Conn
	stream quic.Stream
	local  net.Addr
	remote net.Addr
}

func (c *interConn) Read(b []byte) (int, error) {
	return c.stream.Read(b)
}

func (c *interConn) Write(b []byte) (int, error) {
	return c.stream.Write(b)
}

func (c *interConn) Close() error {
	return c.stream.Close()
}

func (c *interConn) LocalAddr() net.Addr {
	return c.local
}

func (c *interConn) RemoteAddr() net.Addr {
	return c.remote
}

func (c *interConn) SetDeadline(t time.Time) error {
	return c.stream.SetDeadline(t)
}

func (c *interConn) SetReadDeadline(t time.Time) error {
	return c.stream.SetReadDeadline(t)
}

func (c *interConn) SetWriteDeadline(t time.Time) error {
	return c.stream.SetWriteDeadline(t)
}