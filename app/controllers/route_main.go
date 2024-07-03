package controllers

import (
	"fmt"
	// "log"
	"net/http"
	// "strconv"
	// "todo_app/app/models"
	// "path/filepath"
)

func top (w http.ResponseWriter, r *http.Request) {
	generateHTML(w, nil, "layout", "top")
}

func submit_topic (w http.ResponseWriter, r *http.Request) {
	// len := r.ContentLength
	// body := make([]byte, len)
	// r.Body.Read(body)
	// fmt.Fprintln(w, string(body))
	http.Redirect(w, r, "/", 301)
}
