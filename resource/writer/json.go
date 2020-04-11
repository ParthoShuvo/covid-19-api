package writer

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type jsonWriter struct{}

func (w jsonWriter) Write(rw http.ResponseWriter, data interface{}) {
	json, err := json.Marshal(&data)
	if err != nil {
		log.Fatal("json response failed " + err.Error())
		http.Error(rw, "500 - "+err.Error(), http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json; charset=utf8")
	rw.Write(json)
	log.Info("sucessfully responds json data")
}

func newJSONWriter() Writer {
	var jw jsonWriter
	return jw
}
