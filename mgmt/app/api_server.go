package app

import (
	"encoding/json"
	"net/http"
	"time"
)

func sendJSON(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(HookResponse{Code: code, Message: message})
}

func NewHttpServer(listenAddr string) *http.Server {
	if listenAddr == "" {
		listenAddr = ":8080"
	}

	mux := http.NewServeMux()

	// === ALL EVENTS ===
	mux.HandleFunc("/on_connect", handleSrtHook)
	mux.HandleFunc("/check_publish", handleSrtHook)
	mux.HandleFunc("/check_play", handleSrtHook)
	mux.HandleFunc("/check_close", handleSrtHook)
	mux.HandleFunc("/on_play", handleSrtHook)
	mux.HandleFunc("/on_publish", handleSrtHook)
	mux.HandleFunc("/on_unpublish", handleSrtHook)
	mux.HandleFunc("/on_stop", handleSrtHook)
	mux.HandleFunc("/on_dvr", handleSrtHook)
	mux.HandleFunc("/on_hls", handleSrtHook)

	server := &http.Server{
		Addr:         listenAddr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	return server
}
