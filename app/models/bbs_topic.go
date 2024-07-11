package models

import (
	// "fmt"
	// "log"
	"net/http"
	"strings"
	// "strconv"
	// "todo_app/app/models"
	// "path/filepath"
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

var Db *sql.DB
// TODO:configに接続情報をまとめる
var connStr = "user=yudai.kudo dbname=bbs_development sslmode=disable"

type Topic struct {
	Title string
	Description string
	Category string
	FormattedCreatedAt string
}

func GetTopics (w http.ResponseWriter, r *http.Request) (Topics []Topic, err error) {
	
	Db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer Db.Close()

	cmd := `SELECT topic_title, topic_description, topic_category, created_at FROM bbs_topics ORDER BY created_at DESC LIMIT 5`
	rows, err := Db.Query(cmd)
	if err != nil {
		log.Fatalln(err)
	}
	for rows.Next() {
		var topic Topic
        var createdAt time.Time

		weekday := strings.NewReplacer(
			"Sun", "日",
			"Mon", "月",
			"Tue", "火",
			"Wed", "水",
			"Thu", "木",
			"Fri", "金",
			"Sat", "土",
		)

		err = rows.Scan(
			&topic.Title,
			&topic.Description,
			&topic.Category,
			&createdAt,
		)
		if err != nil {
			log.Fatalln(err)
		}
		formattedDate := weekday.Replace(createdAt.Format("2006年1月2日(Mon) 15:04:05"))
		topic.FormattedCreatedAt = formattedDate

		Topics = append(Topics, topic)
	}

	rows.Close()

	return Topics, err

}