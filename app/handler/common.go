package handler

import (
	"encoding/json"
	myerror "github.com/kakugirai/moetify/app/error"
	"log"
	"net/http"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Content interface{} `json:"content"`
}

func RespondWithError(w http.ResponseWriter, err error, payload interface{}) {
	switch e := err.(type) {
	case myerror.Error:
		log.Printf("HTTP %d - %s", e.Status(), e)
		resp, _ := json.Marshal(Response{Code: e.Status(), Message: e.Error(), Content: payload})
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(resp)
	default:
		RespondWithJSON(w, http.StatusInternalServerError, payload)
	}
}

func RespondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	resp, _ := json.Marshal(Response{Code: status, Message: http.StatusText(status), Content: payload})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(resp)
}
