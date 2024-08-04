package controllers

import (
	"fmt"
	//"log"
	"html/template"
	"net/http"
	//"strconv"
	"bbs-development/app/models"
	"bbs-development/config"
	// "context"
)

	type contextKey string
	const sessionKey contextKey = "Session"

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
	http.HandleFunc("/submit_topic", PostTopic)
	http.HandleFunc("/topics/{id}", GetTopic)
	http.HandleFunc("/submit_reply/{id}", PostReply)
	http.HandleFunc("/topics", GetSearchTopicPage)
	http.HandleFunc("/search_topics", SearchTopic)
	http.HandleFunc("/signup", Signup)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/authenticate", Authenticate)
	http.HandleFunc("/mypage", ShowMypage)
	http.HandleFunc("/update_topic", UpdateTopic)
	http.HandleFunc("/delete_topic", DeleteTopic)

	return http.ListenAndServe(":" + config.Config.Port, nil)
}

// TODO : ミドルウェア周りの実装を後々に行う(リクエストのコンテキスト周りの実装が上手くいかなかった）

// func Middleware (next http.HandlerFunc) http.HandlerFunc {
//     return func(w http.ResponseWriter, r *http.Request) {
        
// 		session, _ := session(r)
// 		if &session != nil {
// 			next.ServeHTTP(w, r)
// 		} else {
// 			generateHTML(w, nil, "layout", "login")
// 		}
//     }
// }

func session (r *http.Request) (sess models.Session, err error) {
	cookie , err := r.Cookie("_cookie")

	if err == nil {
		sess = models.Session{UUID: cookie.Value}
		
		if ok, _ := sess.CheckSession(); !ok {
			err = fmt.Errorf("Invalid Session")
		}
	}
	return sess,err
}
