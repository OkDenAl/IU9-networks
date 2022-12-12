package sshWeb

import (
	"bytes"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/ssh"
	"lab10_my/entity"
	"log"
	"net/http"
	"time"
)

func Register(mux *httprouter.Router, cfg *ssh.ClientConfig, upg *websocket.Upgrader) {
	api := provider{
		cfg: cfg,
		upg: upg,
	}

	mux.GET("/ssh", api.handleSSH)
}

type provider struct {
	cfg *ssh.ClientConfig
	upg *websocket.Upgrader
}

func (pr *provider) handleSSH(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	conn, err := pr.upg.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	dial, err := ssh.Dial("tcp", "151.248.113.144:443", pr.cfg)
	if err != nil {
		log.Println(err)
		return
	}
	defer dial.Close()
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			session, err := dial.NewSession()
			if err != nil {
				log.Println(err)
				return
			}

			var outBuff, errBuff bytes.Buffer
			session.Stdout = &outBuff
			session.Stderr = &errBuff
			err = session.Run("cat achtung.txt")
			session.Close()

			message := entity.Response{
				Out: outBuff.String(),
				Err: errBuff.String(),
			}

			if err = conn.WriteJSON(message); err != nil {
				log.Println(err)
				return
			}
		}
	}
}
