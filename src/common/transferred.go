package common

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"github.com/bytedance/gopkg/util/gopool"
)

var translist = make(map[uint16]*TransModel)
var m sync.Mutex

type TransModel struct {
	use       uint64
	port      uint16
	listener  net.Listener
	transfunc func(int) uint
	initFunc  func()
}

func (model *TransModel) Close() {
	model.listener.Close()
}

func CloseTrans(port uint16) error {
	model, ok := translist[port]
	if !ok {
		return fmt.Errorf("关闭异常,未找到端口:%d", port)
	}
	model.use = 0
	model.listener.Close()
	model.listener = nil
	return nil
}

func tcp_transferred(value uint64, sourcePort, destinationAddress string, action func(*TransModel)) error {
	// sourcePort := ":8085"
	// destinationAddress := "127.0.0.1:80"

	port, _, err := GetPort(sourcePort)
	if err != nil {
		return err
	}
	listener, err := net.Listen("tcp", sourcePort)
	if err != nil {
		log.Fatalf("Error listening on port %s: %v", sourcePort, err)
		return err
	}
	defer listener.Close()
	log.Printf("Listening on port %s...\n", sourcePort)
	translist[port] = new(TransModel)
	translist[port].use = value
	translist[port].port = port
	translist[port].listener = listener
	action(translist[port])
	translist[port].initFunc()
	for {
		clientConn, err := listener.Accept()
		if err != nil {
			// log.Printf("Error accepting connection: %v", err)
			clientConn = nil
			return err
		}
		gopool.Go(func() {
			handleConnection(port, clientConn, destinationAddress)
		})
	}
}

func handleConnection(port uint16, clientConn net.Conn, destinationAddress string) {
	defer clientConn.Close()

	destConn, err := net.Dial("tcp", destinationAddress)
	if err != nil {
		// log.Printf("Error connecting to destination: %v", err)
		return
	}
	defer destConn.Close()

	go copyAndCount(port, destConn, clientConn)
	copyAndCount(port, clientConn, destConn)
}

func copyAndCount(port uint16, dst io.Writer, src io.Reader) {
	countedReader := &countingReader{Reader: src, model: translist[port]}
	if _, err := io.Copy(dst, countedReader); err != nil {
		// log.Printf("Error during copy: %v", err)
	}
}

type countingReader struct {
	io.Reader
	model *TransModel
}

func (r *countingReader) Read(p []byte) (int, error) {
	n, err := r.Reader.Read(p)
	if n > 0 {
		m.Lock()
		defer m.Unlock()
		// log.Printf("len: %d", n)
		r.model.use += uint64(n)

		ret := r.model.transfunc(n)
		if ret == 0 { //关闭
			r.model.listener.Close()
			r.model.listener = nil
		} else if ret == 1 { //重新统计
			r.model.use = 0
		}
	}
	return n, err
}
