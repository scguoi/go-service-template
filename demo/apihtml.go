package demo

import (
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

func APIProto(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/html")
	fileBytes, _ := os.ReadFile("demo/index.html")
	_, err := w.Write(fileBytes)
	if err != nil {
		log.Errorf("write response failed. %v", err)
	}
}
