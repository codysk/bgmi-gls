package ts

import (
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"moe.two.bgmi-gls/pkg/service"
	"net/http"
)

func VideoFragment(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()
	path := mux.Vars(r)["path"]
	startAt := vars.Get("startAt")
	duration := vars.Get("duration")

	if path == "" ||
		startAt == "" ||
		duration == "" {
		_, _ = w.Write([]byte("Invalid Param"))
		return
	}

	err := service.TSFragment(path, startAt, duration, w)
	if err != nil {
		log.Error(err)
	}
}
