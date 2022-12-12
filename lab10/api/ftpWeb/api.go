package ftpWeb

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/jlaffaye/ftp"
	"github.com/julienschmidt/httprouter"
	"lab10_my/entity"
	"log"
	"net/http"
	"time"
)

func Register(mux *httprouter.Router, upg *websocket.Upgrader, client *ftp.ServerConn) {
	api := provider{
		upg:    upg,
		client: client,
	}
	mux.GET("/ftp", api.handleFTP)
}

type provider struct {
	upg    *websocket.Upgrader
	client *ftp.ServerConn
}

func (pr *provider) handleFTP(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	upgrade, err := pr.upg.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer upgrade.Close()
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			f, err := pr.client.Retr("achtung.txt")
			defer f.Close()
			var out, errorStr string
			if err != nil {
				errorStr = "no file"
			}
			if err == nil {
				buf := make([]byte, 10)
				for {
					bytesRead, _ := f.Read(buf)
					if bytesRead == 0 {
						break
					}
					out += string(buf)
				}
			}
			resp := &entity.Response{
				Out: out,
				Err: errorStr,
			}
			jsonResp, err := json.Marshal(resp)
			if err != nil {
				log.Println(err)
				return
			}
			err = upgrade.WriteMessage(1, jsonResp)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}
