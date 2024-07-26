package controllers

import (
	"bbs-development/app/models"
	"fmt"
	//"fmt"
	"log"
	"net/http"
)

func Signup (w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {

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
			generateHTML(w, data, "layout", "signup","public_navbar")
		} else {
			generateHTML(w, nil, "layout", "signup","public_navbar")
		}
	} else if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		user := models.User {
			Name: r.PostFormValue("name"),
			Email: r.PostFormValue("email"),
			PassWord: r.PostFormValue("password"),
		}
		if err := user.CreateUser(); err != nil {
			log.Println(err)
		}

		http.Redirect(w, r, "/", 302)
	}
}

func Login (w http.ResponseWriter, r *http.Request) {
	generateHTML(w, nil, "layout", "login")

}

func Authenticate (w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	user, err := models.GetUserByEmail(r.PostFormValue("email"))
	if err != nil {
		log.Fatalln(err)
		http.Redirect(w,r, "/login", 302)
	}
	if user.PassWord == models.Encrypt(r.PostFormValue("password")) {
		session, err := user.CreateSession()
		if err != nil {
			log.Println(err)
		}

		cookie := http.Cookie{
			Name: "_cookie",
			Value: session.UUID,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)

		fmt.Println("認証完了")
		http.Redirect(w,r, "/mypage", 302)
	} else {
		http.Redirect(w,r, "/login", 302)
	}
}
