package server

import (
	"log"
	"net/http"
)

var InnerHTML []byte

func Run() error {
	server := setupServer()

	log.Printf("draftgen server started at: %v", server.Addr)
	return server.ListenAndServe()
}

// setupServer specifies endpoints and corresponding handlers
func setupServer() *http.Server {
	router := http.NewServeMux()
	router.HandleFunc("/", serviceHTMLHandlerFn)
	router.HandleFunc("/generate", generateDraftsHandlerFn)

	server := &http.Server{
		Addr:    ":" + config.Server.Port,
		Handler: router,
	}

	return server
}
