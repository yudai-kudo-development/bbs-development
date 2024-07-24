package models

import (
	"fmt"
	// "log"
	"net/http"
	// "todo_app/app/models"
	"database/sql"
	"log"
	"time"
	//"strings"

	_ "github.com/lib/pq"
)

type Reply struct {
    ID        int
    TopicID   int
    Name      string
    Content   string
    CreatedAt string
}

func PostReply (id int , r *http.Request) (err1 error, err2 error) {
	Db, err1 := sql.Open("postgres", connStr)
	if err1 != nil {
		log.Fatal(err1)
	}
	defer Db.Close()
	
		bbs_topic_id := id
		title := r.FormValue("name")
		description := r.FormValue("content")
	
		insert, err2 := Db.Prepare("INSERT INTO bbs_replies(bbs_topic_id, reply_name, reply_content) VALUES ($1,$2,$3)")
		if err2 != nil {
			fmt.Println(err2)
		}
	
		insert.Exec(bbs_topic_id, title, description)

	return err1,err2
}

func GetReplies (id int) (Replies []Reply, err error) {
	
	var reply Reply
	var createdAt time.Time

	Db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer Db.Close()

	cmd := `SELECT id, reply_name, reply_content, created_at FROM bbs_replies WHERE bbs_topic_id = $1`
	err = Db.QueryRow(cmd, id).Scan(
		&reply.ID,
		&reply.TopicID,
		&reply.Name,
		&reply.Content,
		&createdAt,
	)
	if err != nil {
		log.Fatalln(err)
	}

		Replies = append(Replies, reply)

	return Replies, err
}
