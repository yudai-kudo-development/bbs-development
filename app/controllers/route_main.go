package controllers

import (
	// "fmt"
	// "log"
	"net/http"
	// "strconv"
	"bbs-development/app/models"
	// "path/filepath"
	"database/sql"
	"log"
	"fmt"
	
	_ "github.com/lib/pq"
)

var Db *sql.DB
// TODO:configに接続情報をまとめる
var connStr = "user=yudai.kudo dbname=bbs_development sslmode=disable"

func top (w http.ResponseWriter, r *http.Request) () {
	fmt.Println("Received a request")

	topics, err := models.GetTopics(w,r)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(topics)
	generateHTML(w, topics, "layout", "top")
}

func submit_topic (w http.ResponseWriter, r *http.Request) {

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

		insert, err := Db.Prepare("INSERT INTO bbs_topics(topic_title, topic_description, topic_category) VALUES ($1,$2,$3)")
		if err != nil {
			fmt.Println(err)
		}

		insert.Exec(title, description, category)

	}

	http.Redirect(w, r, "/", 302)

}
