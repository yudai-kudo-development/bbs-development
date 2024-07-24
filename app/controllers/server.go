package controllers

import (
	"fmt"
	//"log"
	"html/template"
	"net/http"
	//"strconv"
	"bbs-development/app/models"
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
	http.HandleFunc("/login", Login)
	http.HandleFunc("/authenticate", Authenticate)
	http.HandleFunc("/mypage", ShowMypage)

	return http.ListenAndServe(":" + config.Config.Port, nil)
}

func session (r *http.Request) (sess models.Session, err error) {
	cookie , err := r.Cookie("_cookie")
	fmt.Println(cookie)

	if err == nil {
		fmt.Println("セッションチェックに入る")
		sess = models.Session{UUID: cookie.Value}
		fmt.Println(sess)
		
		if ok, _ := sess.CheckSession(); !ok {
			fmt.Println("セッションチェック失敗")
			err = fmt.Errorf("Invalid Session")
		}
	}
	return sess,err
}
