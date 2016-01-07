package endpoints

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/gateway/netswitches"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const (
	CMD_ON  = "on"
	CMD_OFF = "off"
)

type HttpServer struct {
	nss *netswitches.NetSwitches
}

func NewHttpServer(netSwitches *netswitches.NetSwitches) (h *HttpServer) {
	h = &HttpServer{
		nss: netSwitches,
	}
	http.HandleFunc("/machines/", h.Handle)
	return
}

func (h *HttpServer) Run() {
	if err := http.ListenAndServe(":7070", nil); err != nil {
		h.nss.Save()
	}
}

func (h *HttpServer) handle(w http.ResponseWriter, r *http.Request) error {
	tmp := strings.Split(r.URL.Path, "/")
	idStr := tmp[len(tmp)-2]
	cmdStr := tmp[len(tmp)-1]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return fmt.Errorf("parse id: %v", err)
	}
	log.Printf("id: %v", id)
	log.Printf("cmd: %v", cmdStr)

	switch cmdStr {
	case CMD_ON, CMD_OFF:
		return h.nss.SetOn(id, cmdStr == CMD_ON)
	default:
		return fmt.Errorf("unknown command '%v'", cmdStr)
	}
	return nil
}

func (h *HttpServer) Handle(w http.ResponseWriter, r *http.Request) {
	if err := h.handle(w, r); err == nil {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error: %v", err)
		log.Printf("run command: %v", err)
	}
}
