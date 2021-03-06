package proxy

import (
	"testing"
)

func TestImplementTCPServer(t *testing.T) {
	var _ Proxy = new(TCPServer)
	var _ Server = new(TCPServer)
}

func TestNew(t *testing.T) {
	s, err := NewTCPServer("127.0.0.1:1081")
	if err != nil {
		t.Error(err)
	}
	//s.Close()
	//select {}
	s.SetServer("127.0.0.1:1082")
	//s.Close()
	select {}
}

func TestDefer(t *testing.T) {
	defer t.Log("main defer")
	s := make(chan bool)
	go func() {
		defer t.Log("defer")
		t.Log("before defer")
		s <- false
	}()

	t.Log("main")
	<-s
}
