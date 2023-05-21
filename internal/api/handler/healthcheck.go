package handler

import (
	"encoding/json"
	"net/http"
)

func HealthCheck(response http.ResponseWriter, request *http.Request) {
	json.NewEncoder(response).Encode(map[string]bool{"ok": true})
}
