package demo

import (
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

func APIProto(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/html")
	fileBytes, _ := os.ReadFile("apiproto/index.html")
	w.Write(fileBytes)
}
