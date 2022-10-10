package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

// 学習用の簡易サーバーなので今回はポートを7777で固定
var port = "7777"

func main() {
	// リクエスト先作成
	service := fmt.Sprintf("%v:%v", ipAddr(os.Args[1]), port)

	// ソケットの作成とIP:portに紐付け
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err, "conn")

	// サーバ側に接続
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err, "conn")

	method := os.Args[2]
	path := os.Args[3]
	// ソケットにデータの書き込み
	reqH := fmt.Sprintf("%v %v HTTP/1.0\r\n", method, path)
	// reqH := fmt.Sprintf("%v %v HTTP/1.0\r\n\r\n", method, path)

	io.WriteString(conn, reqH)
	if method == "POST" || method == "PUT" {
		reqB := os.Args[4]
		buf := []byte(reqB)
		io.WriteString(conn, fmt.Sprintf("Content-Length: %v\r\n\r\n", len(buf)))
		fmt.Println(reqB)
		io.WriteString(conn, reqB)
	}

	res := make([]byte, 1024)

	// ソケットからデータの読み込み
	len, err := conn.Read(res)
	checkError(err, "conn read")
	fmt.Println("response:", string(res[:len]))

	// 接続の切断
	conn.Close()
}

// IP取得
func ipAddr(url string) string {
	addr, err := net.ResolveIPAddr("ip", url)
    if err != nil {
        fmt.Println("Resolve error ", err.Error())
        os.Exit(1)
    }
	return addr.String()
}


func checkError(err error, msg string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s \n", err.Error())
		fmt.Fprintf(os.Stderr, "message: %s \n", msg)
		os.Exit(1)
	}
}