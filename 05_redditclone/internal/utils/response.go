package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func JSONSend(w http.ResponseWriter, data any, code int) {
	respBytes, err := json.Marshal(data)
	if err != nil {
		http.Error(w, `{error: "response serialization failed"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_, err = w.Write(respBytes)
	if err != nil {
		log.Printf("writing data to the connection failed: %v\n", err)
	}
}

func JSONMessage(w http.ResponseWriter, msg string, code int) {
	errObj := map[string]string{"message": msg}
	JSONSend(w, errObj, code)
}
