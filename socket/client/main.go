package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	// ①実行の際に指定したhost:portでbind
	if len(os.Args) != 2 {
        fmt.Fprintf(os.Stderr, "Usage: %s host:port ", os.Args[0])
        os.Exit(1)
    }
    service := os.Args[1]

	// ②ソケットの作成とIP:portに紐付け
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
    checkError(err, "tcpAddr")

	// ③サーバ側に接続
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
    checkError(err, "conn")

	// ④ソケットにデータの書き込み
    _, err = conn.Write([]byte("HEAD / HTTP/1.0\r\n\r\n"))
    checkError(err, "conn write")
    res := make([]byte, 1024)

	// ⑤ソケットからデータの読み込み
    len, err := conn.Read(res)
    checkError(err, "conn read")
    fmt.Println("response:", string(res[:len]))
    // ⑥接続の切断
    conn.Close()
}

func checkError(err error, msg string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s \n", err.Error())
		fmt.Fprintf(os.Stderr, "message: %s \n", msg)
		os.Exit(1)
	}
}