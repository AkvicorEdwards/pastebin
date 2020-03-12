package paste

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"math/rand"
	"net/http"
	"os"
	"pastebin/config"
	"strconv"
	"time"

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
		//r.Body = http.MaxBytesReader(w, r.Body, 32<<20)
		mr, err := r.MultipartReader()
		if err != nil {
			fmt.Println(err)
		}
		
		form := make(map[string]string, 0)

		maxValueBytes := int64(10 << 20)
		type f struct {
			Name string `json:"name"`
			Real string	`json:"real"`
		}
		file := make([]f, 0)
		for {
			part, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			name := part.FormName()
			handleForm := func(name string)  {
				var b bytes.Buffer
				_, err := io.CopyN(&b, part, maxValueBytes)
				if err != nil && err != io.EOF {
					fmt.Sprintln(err)
					fmt.Fprintln(w, err)
				}
				switch name {
				case "highlight":
					form["highlight"] = b.String()
				case "title":
					form["title"] = b.String()
				case "password":
					form["password"] = b.String()
				case "times":
					form["times"] = b.String()
				case "deadline":
					form["deadline"] = b.String()
				case "content":
					form["content"] = b.String()
				}
			}
			if name == "" {
				continue
			}
			if name != "file" {
				handleForm(name)
				continue
			}
			if len(part.FileName()) == 0 {
				continue
			}
			fileName := strconv.FormatInt(time.Now().Unix(), 10)+ "_" + strconv.FormatInt(int64(rand.Int()), 10) + "_" + part.FileName()
			file = append(file, f{
				Name: fileName,
				Real: part.FileName(),
			})

			func(){
				dst, err := os.Create(config.Data.Path.File + fileName)
				if err != nil {
					fmt.Println(err)
					return
				}
				defer dst.Close()
				for {
					buffer := make([]byte, 100000)
					cBytes, errs := part.Read(buffer)

					n, err := dst.Write(buffer[0:cBytes])
					if err != nil {
						fmt.Println(cBytes)
						fmt.Println(n)
						fmt.Println(err)
					}
					if errs == io.EOF {
						break
					}
				}
			}()
		}

		// times
		times := 1000000000
		if len(form["times"]) != 0 {
			times, err =  strconv.Atoi(fmt.Sprint(form["times"]))
			if err != nil {
				_, _ = fmt.Fprint(w, ERROR)
				panic(err)
				return
			}
		}

		// deadline
		deadline := 876000
		if len(form["deadline"]) != 0 {
			deadline, err =  strconv.Atoi(fmt.Sprint(form["deadline"]))
			if err != nil {
				_, _ = fmt.Fprint(w, ERROR)
				panic(err)
				return
			}
		}

		// password
		pwd := ""
		if len(form["password"]) != 0 {
			pwd = fmt.Sprint(form["password"])
		}

		// paste
		paste := bytes.NewBuffer([]byte{})
		jsonEncoder := json.NewEncoder(paste)
		jsonEncoder.SetEscapeHTML(false)
		if err = jsonEncoder.Encode(map[string]interface{}{
			"title": fmt.Sprint(form["title"]),
			"highlight": fmt.Sprint(form["highlight"]),
			"file": file,
			"content": fmt.Sprint(form["content"]),
		}); err != nil {
			_, _ = fmt.Fprint(w, ERROR)
				panic(err)
				return
		}

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
			fmt.Println(err)
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
		if rows.Next() {
			if err = rows.Scan(&password, &paste); err != nil {
				_, _ = fmt.Fprintf(w, `{"status": "%d"}`, ERROR)
				panic(err)
				return
			}
		}else {
			_, _ = fmt.Fprintf(w, `{"status": "%d"}`, NONEXISTENT)
			return
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

		_, _ = db.Exec("update paste set times = times-1 where id = ? limit 1", id)
		_, _ = fmt.Fprintf(w, `{"status": "%d", %s`, ACCEPT, paste[1:])
	}
}

