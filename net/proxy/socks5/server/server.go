package socks5server

import (
	"errors"
	"github.com/Asutorufa/yuhaiin/net/common"
	"github.com/Asutorufa/yuhaiin/net/proxy/socks5/client"
	"net"
	"strconv"
)

// Server <--
type Server struct {
	Username    string
	Password    string
	listener    net.Listener
	udpListener *net.UDPConn
	closed      bool
}

// NewSocks5Server create new socks5 listener
// server: socks5 listener host
// port: socks5 listener port
// username: socks5 server username
// password: socks5 server password
func NewSocks5Server(host, port, username, password string) (s *Server, err error) {
	s = &Server{Username: username, Password: password}
	err = s.Socks5(host, port)
	if err != nil {
		return nil, err
	}
	err = s.UDP(host, port)
	return
}

func (s *Server) UpdateListen(host, port string) (err error) {
	if s.listener.Addr().String() == net.JoinHostPort(host, port) {
		return nil
	}
	if s.listener != nil {
		if err = s.listener.Close(); err != nil {
			return err
		}
	}
	s.listener, err = net.Listen("tcp", net.JoinHostPort(host, port))
	if err != nil {
		return err
	}
	return s.UpdateUDPListenAddr(host, port)
}

func (s *Server) GetListenHost() string {
	return s.listener.Addr().String()
}

// Socks5 <--
func (s *Server) Socks5(host, port string) (err error) {
	s.listener, err = net.Listen("tcp", net.JoinHostPort(host, port))
	if err != nil {
		return err
	}
	go func() {
		for {
			client, err := s.listener.Accept()
			if err != nil {
				if s.closed {
					break
				}
				continue
			}
			_ = client.(*net.TCPConn).SetKeepAlive(true)
			go func() {
				defer client.Close()
				s.handleClientRequest(client)
			}()
		}
	}()
	return
}

// Close close socks5 listener
func (s *Server) Close() error {
	s.closed = true
	s.udpListener.Close()
	return s.listener.Close()
}

func (s *Server) handleClientRequest(client net.Conn) {
	var err error
	b := common.BuffPool.Get().([]byte)
	defer common.BuffPool.Put(b[:cap(b)])

	//socks5 first handshake
	if _, err = client.Read(b[:]); err != nil {
		return
	}

	if b[0] != 0x05 { //只处理Socks5协议
		writeFirstResp(client, 0xff)
		return
	}

	writeFirstResp(client, 0x00)

	if b[1] == 0x01 && b[2] == 0x02 {
		// 对用户名密码进行判断
		if _, err = client.Read(b[:]); err != nil {
			return
		}
		username := b[2 : 2+b[1]]
		password := b[3+b[1] : 3+b[1]+b[2+b[1]]]
		if s.Username != string(username) || s.Password != string(password) {
			writeFirstResp(client, 0x01)
			return
		}
		writeFirstResp(client, 0x00)
	}

	// socks5 second handshake
	_, err = client.Read(b[:])
	if err != nil {
		return
	}

	host, port, _, err := ResolveAddr(b[3:])
	if err != nil {
		return
	}
	//var host, port string
	//switch b[3] {
	//case 0x01: //IP V4
	//	host = net.IPv4(b[4], b[5], b[6], b[7]).String()
	//case 0x03: //domain
	//	host = string(b[5 : n-2]) //b[4] domain's length
	//case 0x04: //IP V6
	//	host = net.IP{b[4], b[5], b[6], b[7], b[8], b[9], b[10], b[11], b[12], b[13], b[14], b[15], b[16], b[17], b[18], b[19]}.String()
	//}
	//port = strconv.Itoa(int(b[n-2])<<8 | int(b[n-1]))

	var server net.Conn
	switch b[1] {
	case 0x01:
		if server, err = common.ForwardTarget(net.JoinHostPort(host, strconv.Itoa(port))); err != nil {
			writeSecondResp(client, 0x04, client.LocalAddr().String())
			return
		}

	case 0x03: // udp
		writeSecondResp(client, 0x00, client.LocalAddr().String())
		for {
			_, err := client.Read(b[:2])
			if err, ok := err.(net.Error); ok && err.Timeout() {
				continue
			}
			return
		}

	case 0x02: // bind request
		fallthrough

	default:
		writeSecondResp(client, 0x07, client.LocalAddr().String())
		return
	}
	defer server.Close()

	writeSecondResp(client, 0x00, client.LocalAddr().String()) // response to connect successful

	// handshake successful
	common.Forward(client, server)
}

func ResolveAddr(raw []byte) (dst string, port, size int, err error) {
	//targetAddrRaw := udpPacket[3:]
	targetAddrRawSize := 1
	switch raw[0] {
	case 0x01:
		dst = net.IP(raw[targetAddrRawSize : targetAddrRawSize+4]).String()
		targetAddrRawSize += 4
	case 0x04:
		if len(raw) < 1+16+2 {
			return "", 0, 0, errors.New("errShortAddrRaw")
		}
		dst = net.IP(raw[1 : 1+16]).String()
		targetAddrRawSize += 16
	case 0x03:
		addrLen := int(raw[1])
		if len(raw) < 1+1+addrLen+2 {
			//return errShortAddrRaw()
		}
		dst = string(raw[1+1 : 1+1+addrLen])
		targetAddrRawSize += 1 + addrLen
	default:
		//s.config.Logger.Printf("udp socks: Failed to get UDP package header: %v", errUnrecognizedAddrType)
		//return errUnrecognizedAddrType
		return "", 0, 0, errors.New("udp socks: Failed to get UDP package header")
	}
	port = (int(raw[targetAddrRawSize]) << 8) | int(raw[targetAddrRawSize+1])
	targetAddrRawSize += 2
	return dst, port, targetAddrRawSize, nil
}

func writeFirstResp(conn net.Conn, errREP byte) {
	_, _ = conn.Write([]byte{0x05, errREP})
}

func writeSecondResp(conn net.Conn, errREP byte, addr string) {
	requestlistenAddr, err := socks5client.ParseAddr(addr)
	if err != nil {
		return
	}
	_, _ = conn.Write(append([]byte{0x05, errREP, 0x00}, requestlistenAddr...))
}
