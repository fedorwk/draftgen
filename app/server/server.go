package server

import (
	"log"
	"net/http"
	"strconv"
)

var InnerHTML []byte

func Run() error {
	err := config.Init()

	config.App.FilenameGenFunc = func(index int, item map[string]string) string {
		return strconv.Itoa(index+1) + ".eml"
	}

	if err != nil {
		return err
	}
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
