package main

import (
	"bufio"
	"errors"
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
		// go handleClient(conn)
		handleClient(conn)

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
	var method, path string
	header := make(map[string]string)

	isFirst := true

	// 一行ずつ処理する
	for {
		line, err := scanner.ReadLine()
		if line == "" {
			break
		}
		if err != nil {
			checkError(err)
		}

		if isFirst {
			isFirst = false
			headerLine := strings.Fields(line)
			header["Method"] = headerLine[0]
			header["Path"] = headerLine[1]
			fmt.Println(method, path)
			continue
		}

		headerFields := strings.SplitN(line, ": ", 2)
        fmt.Printf("%s: %s\n", headerFields[0], headerFields[1])
        header[headerFields[0]] = headerFields[1]
	}

	method, ok := header["Method"]
    if !ok {
        checkError(errors.New("no method found"))
    }

	if method == "POST" || method == "PUT" {
        len, err := strconv.Atoi(header["Content-Length"])
        if err != nil {
            checkError(err)
        }
        buf := make([]byte, len)
        _, err = io.ReadFull(reader, buf)
        if err != nil {
            checkError(err)
        }
        fmt.Println("BODY:", string(buf))
    }
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
