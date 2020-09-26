package m3u8

import (
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"moe.two.bgmi-gls/pkg/service"
	"net/http"
)

func VideoIndex(w http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)["path"]
	if path == "" {
		_, _ = w.Write([]byte("Invalid Param"))
		return
	}
	w.Header().Add("Content-Type", "application/x-mpegURL")
	m3u8, err := service.M3u8VideoIndex(path)
	if err != nil {
		log.Error(err)
	}
	_, _ =w.Write([]byte(m3u8))
}
