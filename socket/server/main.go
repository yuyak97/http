package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	// ①ソケットの作成とIP:portのbind
	service := ":7777"
	// net.ResolveTCPAddr("tcp4", service)の最初の引数は'tcp4', 'tcp6', 'tcp'が設定できます。これはIPv4, IPv6, IPv4orIPv6のどれを使うかを指定しています。
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	// ②接続の待機
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	log.Println("normal socket\nlisten on port", service)
	for {
        // ③接続の受信
        conn, err := listener.Accept()
        if err != nil {
            log.Println(err)
            continue
        }
        // ④ソケットの読み込み
        req := make([]byte, 1024)
        len, err := conn.Read(req)

        log.Println("request:", string(req[:len])) // 部分配列　req[:len]
        // ⑤ソケットの書き込み
        daytime := time.Now().String()
        conn.Write([]byte(daytime))
        // ⑥接続の切断
        conn.Close()
    }
}

func checkError(err error) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
        os.Exit(1)
    }
}