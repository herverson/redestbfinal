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
	li, err := net.Listen("tcp", ":"+port)
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
			url := strings.Fields(ln)[1] // urL
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
		_, _ = fmt.Fprint(conn, "HTTP/1.1 404 Not Found"+CRLF)
		_, _ = fmt.Fprintf(conn, "Content-Length: %d"+CRLF, len(entityBody))
		_, _ = fmt.Fprint(conn, "Content-Type: text/html"+CRLF)
		_, _ = fmt.Fprint(conn, CRLF)
		_, _ = fmt.Fprint(conn, entityBody)
		return
	}

	statusLine := "HTTP/1.1 200 OK"

	file, err := os.Open("." + strings.TrimSpace(fileName)) // For read access.

	if err != nil {
		statusLine = "HTTP/1.1 404 Not Found"
		contentTypeLine := contentType(fileName)
		_ = sendToClient(statusLine, contentTypeLine, conn, file)
		return
	}
	defer file.Close() // make sure to close the file even if we panic.

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
	//.sh	Bourne shell script	application/x-sh
	//.svg	Scalable Vector Graphics (SVG)	image/svg+xml
	//.swf	Small web format (SWF) or Adobe Flash document	application/x-shockwave-flash
	//.tar	Tape Archive (TAR)	application/x-tar
	//.tif
	//.tiff	Tagged Image File Format (TIFF)	image/tiff
	//.ts	Typescript file	application/typescript
	//.ttf	TrueType Font	font/ttf
	//.vsd	Microsoft Visio	application/vnd.visio
	return "application/octet-stream"
}

func sendToClient(statusLine string, contentType string, conn net.Conn, file *os.File) error {

	CRLF := "\r\n"

	if file == nil {
		statusLine += CRLF
		contentType += CRLF
		entityBody := `<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><TITLE>Not Found</TITLE></head><body><strong>Not Found</strong></body></html>`
		_, _ = fmt.Fprint(conn, statusLine)
		_, _ = fmt.Fprintf(conn, "Content-Length: %d\r\n", len(entityBody))
		_, _ = fmt.Fprint(conn, "Content-Type: "+contentType)
		_, _ = fmt.Fprint(conn, CRLF)
		_, _ = fmt.Fprint(conn, entityBody)
		return nil
	}

	statusLine += CRLF
	contentType += CRLF
	_, _ = fmt.Fprint(conn, statusLine)
	_, _ = fmt.Fprint(conn, "Content-Type: "+contentType)
	_, _ = fmt.Fprint(conn, CRLF)
	_, err := io.Copy(conn, file)

	if err != nil {
		return err
	}

	return nil
}
