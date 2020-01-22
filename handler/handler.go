package handler

import (
	"net/http"
	"pastebin/config"
	"pastebin/handler/index"
	"pastebin/handler/paste"
	"regexp"
)

var mux map[string]func(http.ResponseWriter, *http.Request)

// All Prefix
func ParsePrefix() {
	mux = make(map[string]func(http.ResponseWriter, *http.Request))
	mux["/"] = index.Paste
	mux["/api/paste"] = paste.Paste

	mux["/favicon.ico"] = http.FileServer(http.Dir(config.Data.Path.Theme+"img/")).ServeHTTP
}

// Prefix Handler
type MyHandler struct{}

func (*MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h, ok := mux[r.URL.Path]; ok {
		h(w, r)
		return
	}
	if ok, _ := regexp.MatchString("/css/", r.URL.String()); ok {
		http.StripPrefix("/css/", http.FileServer(http.Dir(config.Data.Path.Theme+"css/"))).ServeHTTP(w, r)
	} else if ok, _ := regexp.MatchString("/js/", r.URL.String()); ok {
		http.StripPrefix("/js/", http.FileServer(http.Dir(config.Data.Path.Theme+"js/"))).ServeHTTP(w, r)
	} else if ok, _ := regexp.MatchString("/img/", r.URL.String()); ok {
		http.StripPrefix("/img/", http.FileServer(http.Dir(config.Data.Path.Theme+"img/"))).ServeHTTP(w, r)
	} else {
		mux["/"](w, r)
	}
}

