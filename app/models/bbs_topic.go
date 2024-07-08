package models

import (
	// "fmt"
	// "log"
	"net/http"
	// "strconv"
	// "todo_app/app/models"
	// "path/filepath"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	_ "golang.org/x/text/date"
)

var Db *sql.DB
// TODO:configに接続情報をまとめる
var connStr = "user=yudai.kudo dbname=bbs_development sslmode=disable"

type Topic struct {
	Title string
	Description string
	Category string
}

func GetTopics (w http.ResponseWriter, r *http.Request) (Topics []Topic, err error) {
	
	Db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer Db.Close()

	cmd := `SELECT topic_title, topic_description, topic_category FROM bbs_topics LIMIT 5`
	rows, err := Db.Query(cmd)
	if err != nil {
		log.Fatalln(err)
	}
	for rows.Next() {
		var topic Topic
		err = rows.Scan(
			&topic.Title,
			&topic.Description,
			&topic.Category,
		)
		if err != nil {
			log.Fatalln(err)
		}
		Topics = append(Topics, topic)
	}

	rows.Close()

	return Topics, err

}