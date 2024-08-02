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
	
	Topics, err := models.GetRecentTopics()
	if err != nil {
		fmt.Println(err)
	}

	generateHTML(w, Topics, "layout", "top")

	// session, _ := session(r)

	// if &session != nil {
	// 	User, err := session.GetUserBySession()
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}

	// 	data := map[string]interface{}{
	// 		"Topics": Topics,
	// 		"User": User,
	// 	}

	// 	generateHTML(w, data, "layout", "top")
	// } else {
	// 	generateHTML(w, Topics, "layout", "top")
	// }
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

	Topics, err := models.GetIndividualTopic(id)
	if err != nil {
		fmt.Println(err)
	}

	generateHTML(w, Topics, "layout", "individualtopic")

	// Replies, err := models.GetReplies(id)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// session, _ := session(r)

	// if &session != nil {
	// 	User, err := session.GetUserBySession()
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}

	// 	data := map[string]interface{}{
	// 		"Topics": Topics,
	// 		"Replies":Replies,
	// 		"User": User,
	// 	}

	// 	generateHTML(w, data, "layout", "individualtopic")
	// } else {
	// 	data := map[string]interface{}{
	// 		"Topics": Topics,
	// 		"Replies":Replies,
	// 	}
	// 	generateHTML(w, data, "layout", "individualtopic")
	// }
}

func PostTopic (w http.ResponseWriter, r *http.Request) () {

	session, _ := session(r)

	if &session != nil {
		user, err := session.GetUserBySession()
		if err != nil {
			fmt.Println(err)
		}
		
		err = models.PostTopics(w, r, user.ID)
		if err != nil {
			fmt.Println(err)
		}
		http.Redirect(w, r, "/", 302)
		
	} else {
		var user_id int
		
		err := models.PostTopics(w, r, user_id)
		if err != nil {
			fmt.Println(err)
		}
		http.Redirect(w, r, "/", 302)
	}
}

func PostReply (w http.ResponseWriter, r *http.Request) () {
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

		err = models.PostReply(id, r)
		if err != nil {
			fmt.Println(err)
		}

		redirectURL := fmt.Sprintf("/topics/%d", id)
		http.Redirect(w, r, redirectURL, 302)
	}
}

func GetSearchTopicPage (w http.ResponseWriter, r *http.Request) () {

	session, _ := session(r)

	if &session != nil {
		User, err := session.GetUserBySession()
		if err != nil {
			fmt.Println(err)
		}

		data := map[string]interface{}{
			"Topics": nil,
			"User": User,
		}

		generateHTML(w, data, "layout", "searchtopic")
	} else {
		generateHTML(w, nil, "layout", "searchtopic")
	}
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
	if err != nil {
		http.Redirect(w,r, "/login", 302)
	}

	if &session != nil {

		User, err := session.GetUserBySession()
		if err != nil {
			fmt.Println(err)
		}

		TopicsWithUser, err := models.GetTopicsWithUser(User, User.ID)
		if err != nil {
			fmt.Println(err)
			return
		}

		data := map[string]interface{}{
			"Topics": TopicsWithUser,
			"User": User,
		}

		generateHTML(w, data, "layout", "mypage")
	} else {
		http.Redirect(w, r, "/", 302)
	}
}
