package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	goroutinSocket()
}

func goroutinSocket() {
	// ①ソケットの作成とIP:portのbind
	service := ":7777"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)

	// ②接続の待機
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	log.Println("concurrency socket\nlisten on port", service)
	for {
		// ③接続の受信
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		// 複数のクライアントからの接続を捌く必要があるので並行処理
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	// ⑥接続の切断
	defer conn.Close()
	// ④ソケットの読み込み
	req := make([]byte, 1024)
	len, err := conn.Read(req)
	checkError(err)
	log.Println("request:", string(req[:len])) // 部分配列　req[:len]
	// ⑤ソケットの書き込み
	daytime := time.Now().String()
	conn.Write([]byte(daytime))
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
