package controllers

import "net/http"

func Signup (w http.ResponseWriter, r *http.Request) () {
	generateHTML(w, nil, "layout", "signup","public_navbar")
}
