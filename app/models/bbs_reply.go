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
	"os"
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type BbsReply struct {
    ID           int
    BbsTopicID   int
    ReplyName    string
    ReplyContent string
    CreatedAt    time.Time
		UpdatedAt    time.Time
}

func PostReply (id int , r *http.Request) (err error) {
	err = godotenv.Load(".env")
	if err != nil {
		fmt.Printf("読み込み出来ませんでした: %v", err)
	} 
	Db, err := gorm.Open(postgres.Open(os.Getenv("DB_DSN")), &gorm.Config{})
	if err != nil {
	panic("failed to connect database")
	}

	if r.Method == "POST" {
		reply := BbsReply{
  	  ReplyName: r.FormValue("name"),
  	  ReplyContent: r.FormValue("content"),
			BbsTopicID: id,
		}
	 
		topic_pointer := Db.Create(&reply)
		if topic_pointer.Error != nil {
  	  panic("failed to insert record")
		}
	}

	return err
}

func GetReplies (id int) (Replies []BbsReply, err error) {
	
	var reply BbsReply
	var createdAt time.Time

	Db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer Db.Close()

	cmd := `SELECT reply_name, reply_content, created_at FROM bbs_replies WHERE bbs_topic_id = $1`
	err = Db.QueryRow(cmd, id).Scan(
		&reply.ReplyName,
		&reply.ReplyContent,
		&createdAt,
	)

	if err != nil {
		fmt.Println(err)
		return []BbsReply{}, err
	} else {
		Replies = append(Replies, reply)
		return Replies, err
	}
}
