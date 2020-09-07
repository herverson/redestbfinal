package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	port := "6789"
	li, err := net.Listen("tcp", ":" + port)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer li.Close()

	for {
		conn, err := li.Accept()
		if err != nil {
			log.Println(err.Error())
			continue
		}
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()

	// read request
	request(conn)

}

func request(conn net.Conn) {
	i := 0
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		ln := scanner.Text()
		fmt.Println(ln)
		if i == 0 {
			// request line
			method := strings.Fields(ln)[0] // method
			url := strings.Fields(ln)[1] // urL
			fmt.Println("METHOD", method)
			fmt.Println("URL", url)
			response(url, conn)
		}
		if ln == "" {
			// headers are done
			break
		}
		i++
	}
}

func response(fileName string, conn net.Conn) {
	CRLF := "\r\n"

	d, _ := ioutil.ReadDir("." + fileName)
	if d != nil {
		entityBody := `<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><TITLE>Diretorio</TITLE></head><body><strong>Diretório</strong></body></html>`
		fmt.Fprint(conn, "HTTP/1.1 404 Not Found" + CRLF)
		fmt.Fprintf(conn, "Content-Length: %d" + CRLF, len(entityBody))
		fmt.Fprint(conn, "Content-Type: text/html" + CRLF)
		fmt.Fprint(conn, CRLF)
		fmt.Fprint(conn, entityBody)
		return
	}

	file, err := os.Open("." + strings.TrimSpace(fileName)) // For read access.

	if err != nil {
		entityBody := `<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><TITLE>Not Found</TITLE></head><body><strong>Not Found</strong></body></html>`
		statusLine := "HTTP/1.1 404 Not Found" + CRLF
		fmt.Fprint(conn, statusLine)
		fmt.Fprintf(conn, "Content-Length: %d\r\n", len(entityBody))
		fmt.Fprint(conn, "Content-Type: text/html" + CRLF)
		fmt.Fprint(conn, CRLF)
		fmt.Fprint(conn, entityBody)
		return
	}
	defer file.Close() // make sure to close the file even if we panic.

	fmt.Fprint(conn, "HTTP/1.1 200 OK" + CRLF)
	fmt.Fprint(conn, "Content-Type: text/html" + CRLF)
	fmt.Fprint(conn, CRLF)

	_, err = io.Copy(conn, file)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("bytes sent")
}