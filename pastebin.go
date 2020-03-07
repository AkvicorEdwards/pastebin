package main

import (
	"fmt"
	"net/http"
	"pastebin/config"
	"pastebin/handler"
	"time"
)

func main() {
	config.ParseYaml()
	handler.ParsePrefix()

	server := http.Server{
		Addr:              config.Data.Server.Addr,
		Handler:           &handler.MyHandler{},
		ReadTimeout:       20 * time.Second,
	}

	fmt.Println("ListenAndServe: ", config.Data.Server.Addr)
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
