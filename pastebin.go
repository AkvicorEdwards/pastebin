package main

import (
	"net/http"
	"pastebin/config"
	"pastebin/handler"
	log "pastebin/logger"
	"time"
)

func main() {
	log.Init()
	config.ParseYaml()
	handler.ParsePrefix()

	server := http.Server{
		Addr:              config.Data.Server.Addr,
		Handler:           &handler.MyHandler{},
		ReadTimeout:       20 * time.Second,
	}
	log.Behaviour.Print("ListenAndServe: ", config.Data.Server.Addr)
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
