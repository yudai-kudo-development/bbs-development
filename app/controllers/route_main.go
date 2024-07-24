package controllers

import (
	"bbs-development/app/models"
	"database/sql"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	_ "github.com/lib/pq"
)

var Db *sql.DB

func Top (w http.ResponseWriter, r *http.Request) () {
	Topics, err := models.GetRecentTopics(w,r)
	if err != nil {
		fmt.Println(err)
	}
	generateHTML(w, Topics, "layout", "top")
}

func GetTopic (w http.ResponseWriter, r *http.Request) () {
	idStr := filepath.Base(r.URL.Path)
	if idStr == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	topics, err := models.GetIndividualTopic(id)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Reply: %+v\n", topics)
	generateHTML(w, topics, "layout", "individualtopic")
}

func SunbmitTopic (w http.ResponseWriter, r *http.Request) () {
	_, err := models.PostTopics(w,r)
	if err != nil {
		fmt.Println(err)
	}
	http.Redirect(w, r, "/", 302)
}

func SunbmitReply (w http.ResponseWriter, r *http.Request) () {
	if r.Method == "POST" {
	idStr := filepath.Base(r.URL.Path)
	if idStr == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	_, err = models.PostReply(id, r)
	if err != nil {
		fmt.Println(err)
	}

	redirectURL := fmt.Sprintf("/topics/%d", id)
	http.Redirect(w, r, redirectURL, 302)
	}
}

func GetSearchTopicPage (w http.ResponseWriter, r *http.Request) () {
	data := ""
	generateHTML(w, data, "layout", "searchtopic")
}

func SearchTopic (w http.ResponseWriter, r *http.Request) () {

	if r.Method == "POST" {
		Topics, err := models.GetTopics(w,r)
		if err != nil {
			fmt.Println(err)
		}
		
		fmt.Println(Topics)
		generateHTML(w, Topics, "layout", "searchtopic")
	}
}

func ShowMypage (w http.ResponseWriter, r *http.Request) () {

	session, err := session(r)
	fmt.Println("ShowMypageのsession抜ける")
	if err != nil {
		http.Redirect(w,r, "/login", 302)
	}

		fmt.Println("elseの処理に入る")
		user, err := session.GetUserBySession()
		if err != nil {
			fmt.Println("user返すところでエラー")
			fmt.Println(err)
		}

	generateHTML(w, user, "layout", "mypage")
}


