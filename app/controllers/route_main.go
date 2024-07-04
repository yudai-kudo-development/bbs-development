package controllers

import (
	// "fmt"
	// "log"
	"net/http"
	// "strconv"
	// "todo_app/app/models"
	// "path/filepath"
	"database/sql"
	"log"
	"fmt"

	_ "github.com/lib/pq"
)

func top (w http.ResponseWriter, r *http.Request) {
	generateHTML(w, nil, "layout", "top")
}

var Db *sql.DB

func submit_topic (w http.ResponseWriter, r *http.Request) {

		// TODO:configに接続情報をまとめる
		connStr := "user=yudai.kudo dbname=bbs_development sslmode=disable"

		Db, err := sql.Open("postgres", connStr)
		if err != nil {
			log.Fatal(err)
		}
		defer Db.Close()
	
		fmt.Println("DB接続完了")
	
	if r.Method == "POST" {
		title := r.FormValue("title")
		description := r.FormValue("description")
		category := r.FormValue("category")

		insert, err := Db.Prepare("INSERT INTO bbs_topics(topic_title, topic_description, topic_categpry) VALUES ($1,$2,$3)")
		if err != nil {
			fmt.Println(err)
		}

		insert.Exec(title, description, category)

	}

	http.Redirect(w, r, "/", 302)

}
