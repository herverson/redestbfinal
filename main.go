package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"os"
	"strings"
)

func main() {
	port := "6789"
	// listen na porta "port"
	li, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer li.Close()
	// laço infinito (ou ctrl-c)
	for {
		// aceita a conexão com o client
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

	// ler a request
	request(conn)

}

func request(conn net.Conn) {
	i := 0
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		ln := scanner.Text()
		fmt.Println(ln)
		if i == 0 {
			// linha do request
			url := strings.Fields(ln)[1] // urL
			respond(url, conn)
		}
		if ln == "" {
			//cabeçalhos estão prontos
			break
		}
		i++
	}
}

func respond(fileName string, conn net.Conn) {

	statusLine := "200 OK"

	convertedFilename, _ := url.PathUnescape(fileName)

	file, err := os.Open("." + convertedFilename)
	// verifica a estrutura do arquivo
	fi, _ := file.Stat()

	// caso não existe um arquivo ou é um diretório
	if err != nil || fi.IsDir() {
		statusLine = "404 Not Found"
		_ = sendToClient(statusLine, "text/html", conn, file)
		return
	}
	defer file.Close()

	contentTypeLine := contentType(fileName)
	err = sendToClient(statusLine, contentTypeLine, conn, file)

	if err != nil {
		log.Fatal(err)
	}

}

func contentType(fileName string) string {

	if strings.HasSuffix(fileName, ".html") || strings.HasSuffix(fileName, ".htm") {
		return "text/html"
	}
	if strings.HasSuffix(fileName, ".txt") {
		return "text/plain"
	}
	//image/gif, image/png, image/jpeg, image/bmp, image/webp
	if strings.HasSuffix(fileName, ".gif") {
		return "image/gif"
	}
	if strings.HasSuffix(fileName, ".png") {
		return "image/png"
	}
	if strings.HasSuffix(fileName, ".jpg") {
		return "image/jpg"
	}
	if strings.HasSuffix(fileName, ".jpeg") {
		return "image/jpeg"
	}
	if strings.HasSuffix(fileName, ".bmp") {
		return "image/bmp"
	}
	if strings.HasSuffix(fileName, ".webp") {
		return "image/webp"
	}
	//.pdf	Adobe Portable Document Format (PDF)	application/pdf
	if strings.HasSuffix(fileName, ".pdf") {
		return "application/pdf"
	}

	//.ppt	Microsoft PowerPoint	application/vnd.ms-powerpoint
	if strings.HasSuffix(fileName, ".ppt") {
		return "application/vnd.ms-powerpoint"
	}
	//.rar	RAR archive	application/x-rar-compressed
	if strings.HasSuffix(fileName, ".rar") {
		return "application/x-rar-compressed"
	}
	//.rtf	Rich Text Format (RTF)	application/rtf
	if strings.HasSuffix(fileName, ".rtf") {
		return "application/rtf"
	}
	//.sh	Bourne shell script	application/x-sh
	if strings.HasSuffix(fileName, ".sh") {
		return "application/x-sh"
	}
	//.svg	Scalable Vector Graphics (SVG)	image/svg+xml
	if strings.HasSuffix(fileName, ".svg") {
		return "image/svg+xml"
	}
	//.swf	Small web format (SWF) or Adobe Flash document	application/x-shockwave-flash
	if strings.HasSuffix(fileName, ".swf") {
		return "application/x-shockwave-flash"
	}
	//.tar	Tape Archive (TAR)	application/x-tar
	if strings.HasSuffix(fileName, ".tar") {
		return "application/x-tar"
	}
	//.tif
	//.tiff	Tagged Image File Format (TIFF)	image/tiff
	if strings.HasSuffix(fileName, ".tiff") {
		return "image/tiff"
	}
	//.ts	Typescript file	application/typescript
	if strings.HasSuffix(fileName, ".ts") {
		return "application/typescript"
	}
	//.ttf	TrueType Font	font/ttf
	if strings.HasSuffix(fileName, ".ttf") {
		return "font/ttf"
	}
	//.vsd	Microsoft Visio	application/vnd.visio
	if strings.HasSuffix(fileName, ".vsd") {
		return "application/vnd.visio"
	}
		return "application/octet-stream"
}

func sendToClient(statusLine string, contentType string, conn net.Conn, file *os.File) error {

	CRLF := "\r\n"
	statusLine += CRLF
	contentType += CRLF

	if statusLine == "404 Not Found\r\n" {

		entityBody := `<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><TITLE>Not Found</TITLE></head><body><strong>Not Found</strong></body></html>`
		_, _ = fmt.Fprint(conn, "HTTP/1.0 " + statusLine)
		_, _ = fmt.Fprintf(conn, "Content-Length: %d\r\n", len(entityBody))
		_, _ = fmt.Fprint(conn, "Content-Type: " + contentType)
		_, _ = fmt.Fprint(conn, CRLF)
		_, _ = fmt.Fprint(conn, entityBody)
		return nil
	}

	_, _ = fmt.Fprint(conn, "HTTP/1.0 " + statusLine)
	_, _ = fmt.Fprint(conn, "Content-Type: " + contentType)
	_, _ = fmt.Fprint(conn, CRLF)
	//_, err := io.Copy(conn, file)
	err := sendBytes(conn, file)

	if err != nil {
		return err
	}

	return nil
}

func sendBytes(conn net.Conn, file *os.File) error {
	// make um buffer para armazenar os pedaços que são lidos
	buf := make([]byte, 1024)
	for {
		// ler um pedaço do arquivo
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}

		if n == 0 {
			break
		}
		// escreve o pedaço
		if _, err := conn.Write(buf[:n]); err != nil {
			return err
		}
	}
	return nil
}
