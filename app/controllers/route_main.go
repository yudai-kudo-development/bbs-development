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
	fmt.Println("top関数の処理に入る")
	generateHTML(w, nil, "layout", "top")
}