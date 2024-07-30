package src

import (
	"io"
	"log"
	"net"
	"sync"
)

type transFunc func(uint64, int) uint

var trans map[string]uint64 = make(map[string]uint64)
var m sync.Mutex
var transfunc transFunc

func Transferred(value uint64, sourcePort, destinationAddress string, f transFunc) {
	// sourcePort := ":8085"
	// destinationAddress := "127.0.0.1:57890"
	trans[sourcePort] = value
	transfunc = f

	listener, err := net.Listen("tcp", sourcePort)
	if err != nil {
		log.Fatalf("Error listening on port %s: %v", sourcePort, err)
	}
	defer listener.Close()
	log.Printf("Listening on port %s...\n", sourcePort)

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			// log.Printf("Error accepting connection: %v", err)
			clientConn = nil
			return
		}
		go handleConnection(listener, sourcePort, clientConn, destinationAddress)
	}
}

func handleConnection(listener net.Listener, key string, clientConn net.Conn, destinationAddress string) {
	defer clientConn.Close()

	destConn, err := net.Dial("tcp", destinationAddress)
	if err != nil {
		// log.Printf("Error connecting to destination: %v", err)
		return
	}
	defer destConn.Close()

	go copyAndCount(listener, key, destConn, clientConn)
	copyAndCount(listener, key, clientConn, destConn)
}

func copyAndCount(listener net.Listener, key string, dst io.Writer, src io.Reader) {
	countedReader := &countingReader{Reader: src, key: key, listener: listener}
	if _, err := io.Copy(dst, countedReader); err != nil {
		// log.Printf("Error during copy: %v", err)
	}
}

type countingReader struct {
	io.Reader
	key      string
	listener net.Listener
}

func (r *countingReader) Read(p []byte) (int, error) {
	n, err := r.Reader.Read(p)
	if n > 0 {
		m.Lock()
		defer m.Unlock()
		trans[r.key] += uint64(n)

		ret := transfunc(trans[r.key], n)
		if ret == 0 { //关闭
			r.listener.Close()
			r.listener = nil
			trans[r.key] = 0
		} else if ret == 1 { //重新统计
			trans[r.key] = 0
		}
	}
	return n, err
}
