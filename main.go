package main

import (
	"fmt"
	"golang.org/x/crypto/acme/autocert"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

func main() {
	host := os.Getenv("HOST")
	var ln net.Listener
	if host == "" || host == "localhost" || host == "127.0.0.1" {
		ln, _ = net.Listen("tcp", ":80")
	} else {
		ln = autocert.NewListener(host)
	}

	log.Fatalln(http.Serve(ln, new(serve)))
}

// serve - どんなリクエストでもログにはくサーバ
type serve struct{}

func (s *serve) ServeHTTP(_ http.ResponseWriter, request *http.Request) {
	file, _ := os.OpenFile(fmt.Sprintf("logs/request-%s.log", time.Now().Format("20060102")), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	defer file.Close()
	l := log.New(file, "", log.LstdFlags|log.Lmicroseconds)

	defer request.Body.Close()
	body, err := ioutil.ReadAll(request.Body)
	l.Printf("request: %+v, body: %s, readBodyErr: %v\n", request, body, err)
}
