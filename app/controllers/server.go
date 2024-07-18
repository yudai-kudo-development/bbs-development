package controllers

import (
	"fmt"
	//"log"
	"html/template"
	"net/http"
	//"strconv"
	//"todo_app/app/models"
	"bbs-development/config"
)

func generateHTML (w http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("app/views/templates/%s.html", file))
	}

	templates := template.Must(template.ParseFiles(files...))
	templates.ExecuteTemplate(w, "layout", data)
}

func StartMainServer() error {
	files := http.FileServer(http.Dir(config.Config.Static))
	http.Handle("/static/", http.StripPrefix("/static/", files))

	http.HandleFunc("/", Top)
	http.HandleFunc("/submit_topic", SunbmitTopic)
	http.HandleFunc("/topics/{id}", GetTopic)
	http.HandleFunc("/submit_reply/{id}", SunbmitReply)
	http.HandleFunc("/topics", GetSearchTopicPage)
	http.HandleFunc("/search_topics", SearchTopic)
	http.HandleFunc("/signup", Signup)

	return http.ListenAndServe(":" + config.Config.Port, nil)
}
