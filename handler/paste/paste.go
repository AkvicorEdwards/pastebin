package paste

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"pastebin/config"
	log "pastebin/logger"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

const (
	ACCEPT = iota
	ERROR
	NONEXISTENT
	WRONGPASSWORD
)

func Paste(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s?charset=%s",
			config.Data.Mysql.User, config.Data.Mysql.Password, config.Data.Mysql.Database, config.Data.Mysql.Charset))
		if err != nil {
			_, _ = fmt.Fprint(w, ERROR)
			panic(err)
			return
		}
		defer func() {
			if err = db.Close(); err != nil {
				panic(err)
			}
		}()

		//// TODO
		//// 处理表单文件
		//// 根据字段名获取表单文件
		//formFile, header, err := r.FormFile("uploadfile")
		//if err != nil {
		//	log.Behaviour.Printf("Get form file failed: %s\n", err)
		//	return
		//}
		//defer formFile.Close()
		//// 创建保存文件
		//destFile, err := os.Create("./upload/" + header.Filename)
		//if err != nil {
		//	log.Behaviour.Printf("Create failed: %s\n", err)
		//	return
		//}
		//defer destFile.Close()
		//
		//// 读取表单文件，写入保存文件
		//_, err = io.Copy(destFile, formFile)
		//if err != nil {
		//	log.Behaviour.Printf("Write file failed: %s\n", err)
		//	return
		//}
		//// 处理文件结束

		if err = r.ParseForm(); err != nil {
			_, _ = fmt.Fprint(w, ERROR)
			panic(err)
			return
		}

		// times
		times := 1000000000
		if len(r.PostForm["times"][0]) != 0 {
			times, err =  strconv.Atoi(fmt.Sprint(r.PostForm["times"][0]))
			if err != nil {
				_, _ = fmt.Fprint(w, ERROR)
				panic(err)
				return
			}
		}

		// deadline
		deadline := 876000
		if len(r.PostForm["deadline"][0]) != 0 {
			deadline, err =  strconv.Atoi(fmt.Sprint(r.PostForm["deadline"][0]))
			if err != nil {
				_, _ = fmt.Fprint(w, ERROR)
				panic(err)
				return
			}
		}

		// password
		pwd := ""
		if len(r.PostForm["password"][0]) != 0 {
			pwd = fmt.Sprint(r.PostForm["password"][0])
		}

		// paste
		paste := bytes.NewBuffer([]byte{})
		jsonEncoder := json.NewEncoder(paste)
		jsonEncoder.SetEscapeHTML(false)
		if err = jsonEncoder.Encode(map[string]string{
			"title": fmt.Sprint(r.PostForm["title"][0]),
			"highlight": fmt.Sprint(r.PostForm["highlight"][0]),
			"content": fmt.Sprint(r.PostForm["content"][0]),
		}); err != nil {
			_, _ = fmt.Fprint(w, ERROR)
				panic(err)
				return
		}
		//fmt.Fprintln(w, "password: [", pwd, "] times: [", times, "] deadline: [", deadline,
		//	"] \npaste:\n", paste.String())

		result, err := db.Exec("insert into paste (pwd, times, deadline, paste) " +
			"values(?, ?, date_add(now(), interval ? hour), ?)",
			pwd, times, deadline, paste.String())
		if err != nil {
			_, _ = fmt.Fprint(w, ERROR)
			panic(err)
			return
		}
		id, err := result.LastInsertId()

		t, _ := template.ParseFiles(config.Data.Path.Theme + "paste_success.html")
		phrase := map[string]interface{}{
			"lang": "zh-cn",
			"title": "Akvicor's PasteBin",
			"url": fmt.Sprintf("http://%s/", r.Host),
			"id": id,
		}

		if err := t.Execute(w, phrase); err != nil {
			log.Behaviour.Println(err)
			_, _ = fmt.Fprintf(w, "%v", "Error")
		}
	}else if r.Method == "GET" {
		if err := r.ParseForm(); err != nil {
			_, _ = fmt.Fprintf(w, `{"status": "%d"}`, ERROR)
			return
		}

		id := -1
		pwd := ""
		for k, v := range r.Form {
			if k == "id" {
				var err error
				id, err =  strconv.Atoi(fmt.Sprint(v[0]))
				if err != nil {
					_, _ = fmt.Fprintf(w, `{"status": "%d"}`, ERROR)
					return
				}
			}else if k == "pwd" {
				pwd = v[0]
			}
		}
		if id == -1 {
			_, _ = fmt.Fprintf(w, `{"status": "%d"}`, NONEXISTENT)
			return
		}

		db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s?charset=%s",
			config.Data.Mysql.User, config.Data.Mysql.Password, config.Data.Mysql.Database, config.Data.Mysql.Charset))
		if err != nil {
			_, _ = fmt.Fprintf(w, `{"status": "%d"}`, ERROR)
			panic(err)
			return
		}
		defer func() {
			if err = db.Close(); err != nil {
				_, _ = fmt.Fprintf(w, `{"status": "%d"}`, ERROR)
				panic(err)
				return
			}
		}()
		rows, err := db.Query("select pwd, paste from paste where id = ? and deadline > now() and times > 0 limit 1", id)
		if err != nil {
			_, _ = fmt.Fprintf(w, `{"status": "%d"}`, ERROR)
			panic(err)
			return
		}
		paste := ""
		password := ""
		for rows.Next() {
			if err = rows.Scan(&password, &paste); err != nil {
				_, _ = fmt.Fprintf(w, `{"status": "%d"}`, ERROR)
				panic(err)
				return
			}
		}
		//fmt.Println(paste)
		if len(paste) == 0 {
			_, _ = fmt.Fprintf(w, `{"status": "%d"}`, NONEXISTENT)
			return
		}
		if pwd != password {
			_, _ = fmt.Fprintf(w, `{"status": "%d"}`, WRONGPASSWORD)
			return
		}

		_, _ = fmt.Fprintf(w, `{"status": "%d", %s`, ACCEPT, paste[1:])
	}
}
