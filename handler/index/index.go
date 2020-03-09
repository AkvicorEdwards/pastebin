package index

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"pastebin/config"
	"strings"
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
			fmt.Println(err)
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
			fmt.Println(err)
			_, _ = fmt.Fprintf(w, "%v", "Error")
		}

	}


	return
}

func Raw(w http.ResponseWriter, r *http.Request) {
	id := strings.Split(fmt.Sprintf("%s", r.URL)[5:], "?")[0]
	que, ok := r.URL.Query()["pwd"]
	pwd := ""
	if ok {
		pwd = que[0]
	}
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s?charset=%s",
		config.Data.Mysql.User, config.Data.Mysql.Password, config.Data.Mysql.Database, config.Data.Mysql.Charset))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err = db.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	rows, err := db.Query("select pwd, paste from paste where id = ? and deadline > now() and times > 0 limit 1", id)
	if err != nil {
		fmt.Println(err)
		return
	}

	paste := ""
	password := ""
	if rows.Next() {
		if err = rows.Scan(&password, &paste); err != nil {
			fmt.Println(err)
			return
		}
	}else {
		return
	}
	if len(password)!=0 && pwd != password {
		_, _ = fmt.Fprint(w, "Need Password")
		return
	}
	_, _ = db.Exec("update paste set times = times-1 where id = ? limit 1", id)
	content := struct {
		Content string `json:"content"`
	}{}
	err = json.Unmarshal([]byte(paste), &content)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, _ = fmt.Fprint(w, content.Content)
	return
}