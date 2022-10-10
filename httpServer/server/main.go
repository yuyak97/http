package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/textproto"
	"os"
	"path/filepath"
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
	method, path, body := scanRequest(scanner, reader, conn)
	route(conn, method, path, body)

	fmt.Println("<<< end")
}

//  HTTP リクエストをパース
func scanRequest(scanner *textproto.Reader, reader *bufio.Reader, conn net.Conn) (method, path, body string) {
	header := make(map[string]string)
	isFirst := true

	// 一行ずつ処理する
	for {
		line, err := scanner.ReadLine()
		// 返り値が空文字であれば空行と判断する
		if line == "" {
			break
		}
		
		checkError(err)
	
		if isFirst {
			isFirst = false
			headerLine := strings.Fields(line)
			header["Method"] = headerLine[0]
			header["Path"] = headerLine[1]
			continue
		}

		headerFields := strings.SplitN(line, ": ", 2)
        header[headerFields[0]] = headerFields[1]
	}

	// リクエストを表示
	for k, v := range header {
		fmt.Printf("%v:%v\n", k, v)
	}

	method = header["Method"]
	path = header["Path"]
	
	if method == "POST" || method == "PUT" {
        len, err := strconv.Atoi(header["Content-Length"])
       
        checkError(err)
        
        buf := make([]byte, len)
        _, err = io.ReadFull(reader, buf)
        
        checkError(err)
        
        fmt.Println("BODY:", string(buf))

		body = string(buf)
    }

	return method, path, body
}

// ルーティング
func route(conn net.Conn, method, path, body string) {
	if method == "GET" && path == "/" {
		io.WriteString(conn, "HTTP/1.1 200 OK\r\n")
		io.WriteString(conn, "Content-Type: text/html\r\n")
    	io.WriteString(conn, "\r\n")
    	io.WriteString(conn, "<h1>Hello World!!</h1>")
	} else if method == "GET" && path == "/sample.html" {
		cwd, err := os.Getwd()
		checkError(err)
		p := filepath.Join(cwd, filepath.Clean(path))
		resp, err := ioutil.ReadFile(p)
		io.WriteString(conn, "HTTP/1.1 200 OK\r\n")
		io.WriteString(conn, "Content-Type: text/html\r\n")
    	io.WriteString(conn, "\r\n")
    	io.WriteString(conn, string(resp))
	} else if method == "GET" && path == "/hello" {
		r := `{"message": "hello"}`
		io.WriteString(conn, "HTTP/1.1 200 OK\r\n")
		io.WriteString(conn, "Content-Type: application/json;charset=UTF-8\r\n")
    	io.WriteString(conn, "\r\n")
    	io.WriteString(conn, r)
	} else if method == "POST" && path == "/hello" {
		io.WriteString(conn, "HTTP/1.1 200 OK\r\n")
		io.WriteString(conn, "Content-Type: application/json;charset=UTF-8\r\n")
    	io.WriteString(conn, "\r\n")
    	io.WriteString(conn, body)
	} else {
		io.WriteString(conn, "HTTP/1.1 404 OK\r\n")
		io.WriteString(conn, "Content-Type: text/plain\r\n")
    	io.WriteString(conn, "\r\n")
    	io.WriteString(conn, "404")
	}
}

// エラーチェック
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
