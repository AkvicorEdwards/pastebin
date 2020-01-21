package index

import (
	"fmt"
	"html/template"
	"net/http"
	"pastebin/config"
	log "pastebin/logger"
)

// Index page
func Index(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles(config.Data.Path.Theme + "index.html")
	phrase := map[string]interface{}{
		"lang": "en",
		"title": "Akvicor's Paste bin",
		"h1": "Title",
		"type": "markdown",
		"code": `
#include <bits/stdc++.h>

using namespace std;

int main() {
	int a, b;
	cin >> a >> b;
	cout << a+b << endl;
	return 0;
}
`,
	}
	if err := t.Execute(w, phrase); err != nil {
		log.Behaviour.Println(err)
		_, _ = fmt.Fprintf(w, "%v", "Error")
	}
	return
}

func Paste(w http.ResponseWriter, r *http.Request) {
	//log.Behaviour.Info(fmt.Sprintf("%s%s", r.Host, r.URL))
	id := fmt.Sprintf("%s", r.URL)[1:]
	url := fmt.Sprintf("http://%s/api/paste?id=%s", r.Host, id)
	urlShort := fmt.Sprintf("http://%s/", r.Host)

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
		url := fmt.Sprintf("http://%s/api/paste?id=%s", r.Host, id)
		urlShort := fmt.Sprintf("http://%s/", r.Host)

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

