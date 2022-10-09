package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/textproto"
	"os"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("start tcp listen...")

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
	fmt.Println(">>> start")

	// ⑥接続の切断
	defer conn.Close()

	reader := bufio.NewReader(conn)
	scanner := textproto.NewReader(reader)
	scanRequest(scanner, reader, conn)

	fmt.Println("<<< end")
}

func scanRequest(scanner *textproto.Reader, reader *bufio.Reader, conn net.Conn) {
	var contentLength int

	// 一行ずつ処理する
	for {
		line, err := scanner.ReadLine()
		if line == "" {
			break
		}
		if err != nil {
			checkError(err)
		}

		if strings.HasPrefix(line, "Content-Length") {
			contentLength, err = strconv.Atoi(strings.TrimSpace(strings.Split(line, ":")[1]))
			if err != nil {
				checkError(err)
			}
		}
		fmt.Println(line)
	}

	// リクエストボディ
	buf := make([]byte, contentLength)
	_, err := io.ReadFull(reader, buf)
	if err != nil {

	}
	fmt.Println("BODY:", string(buf))

	if err != nil {
		checkError(err)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
