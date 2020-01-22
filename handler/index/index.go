package index

import (
	"fmt"
	"html/template"
	"net/http"
	"pastebin/config"
	log "pastebin/logger"
)

func Paste(w http.ResponseWriter, r *http.Request) {
	id := fmt.Sprintf("%s", r.URL)[1:]
	url := fmt.Sprintf("%s://%s/api/paste?id=%s", config.Data.Server.Protocol, r.Host, id)
	urlShort := fmt.Sprintf("%s://%s/", config.Data.Server.Protocol, r.Host)

	if len(fmt.Sprint(r.URL)) == 1 {
		t, _ := template.ParseFiles(config.Data.Path.Theme + "paste.html")
		phrase := map[string]interface{}{
			"lang": "zh-cn",
			"title": "Akvicor's PasteBin",
			"url": url,
			"urlShort": urlShort,
		}
		if err := t.Execute(w, phrase); err != nil {
			log.Behaviour.Println(err)
			_, _ = fmt.Fprintf(w, "%v", "Error")
		}
	}else {
		id := fmt.Sprintf("%s", r.URL)[1:]
		url := fmt.Sprintf("%s://%s/api/paste?id=%s", config.Data.Server.Protocol, r.Host, id)
		urlShort := fmt.Sprintf("%s://%s/", config.Data.Server.Protocol, r.Host)

		t, _ := template.ParseFiles(config.Data.Path.Theme + "index.html")
		phrase := map[string]interface{}{
			"lang": "zh-cn",
			"title": "Akvicor's PasteBin",
			"url": url,
			"urlShort": urlShort,
		}

		if err := t.Execute(w, phrase); err != nil {
			log.Behaviour.Println(err)
			_, _ = fmt.Fprintf(w, "%v", "Error")
		}

	}


	return
}

