package main

import (
	"github.com/gorilla/websocket"
	"github.com/jlaffaye/ftp"
	"github.com/julienschmidt/httprouter"
	"lab10_my/api/ftpWeb"
	"log"
	"net/http"
)

func main() {
	mux := httprouter.New()

	upgrade := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	client, err := ftp.Dial("students.yss.su:21")
	if err != nil {
		log.Fatal(err)
	}
	err = client.Login("ftpiu8", "3Ru7yOTA")
	if err != nil {
		log.Fatal(err)
	}

	ftpWeb.Register(mux, upgrade, client)

	log.Fatal(http.ListenAndServe("localhost:8181", mux))
}
