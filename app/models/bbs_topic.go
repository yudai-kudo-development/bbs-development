package models

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"os"
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type BbsTopic struct {
  	ID               int
  	TopicTitle       string
  	TopicDescription string
		TopicCategory    string
  	BbsUserId        int
		CreatedAt        time.Time
		UpdatedAt        time.Time
		BbsReplies       []BbsReply
}

func GetRecentTopics () (Topics []BbsTopic, err error) {
	
	err = godotenv.Load(".env")
	if err != nil {
		fmt.Printf("読み込み出来ませんでした: %v", err)
	} 
	Db, err := gorm.Open(postgres.Open(os.Getenv("DB_DSN")), &gorm.Config{})
	if err != nil {
	panic("failed to connect database")
	}

	var topic []BbsTopic
	topic_pointer := Db.Order("id desc").Limit(5).Find(&topic)

	if topic_pointer.Error != nil {
		panic("failed to insert record")
	}

	return topic,err
}

func GetIndividualTopic (id int) ([]BbsTopic, error) {

	err = godotenv.Load(".env")
	if err != nil {
		fmt.Printf("読み込み出来ませんでした: %v", err)
	} 
	Db, err := gorm.Open(postgres.Open(os.Getenv("DB_DSN")), &gorm.Config{})
	if err != nil {
	panic("failed to connect database")
	}

	var topic []BbsTopic
	err = Db.Preload("BbsReplies").First(&topic, id).Error
	fmt.Println(topic)
	// bbs_topic := `SELECT id, topic_title, topic_description, topic_category, created_at FROM bbs_topics WHERE id = $1`
	// // err = Db.QueryRow(bbs_topic, id).Scan(&topic.ID, &topic.Topic_title, &topic.Topic_description, &topic.Topic_category)
	// if err != nil {
	// 		return topic, err
	// }

	// // bbs_replies := `SELECT id, bbs_topic_id, reply_name, reply_content, created_at FROM bbs_replies WHERE bbs_topic_id = $1`
	// // rows, err := Db.Query(bbs_replies, id)
	// // if err != nil {
	// // 		return topic, err
	// // }
	// // defer rows.Close()

	// // for rows.Next() {
	// // 		var reply Reply
	// // 		err := rows.Scan(&reply.ID, &reply.TopicID, &reply.Name, &reply.Content, &reply.CreatedAt)
	// // 		if err != nil {
	// // 				return topic, err
	// // 		}
	// // 		topic.Replies = append(topic.Replies, reply)
	// // }

	// // if err := rows.Err(); err != nil {
	// // 	return topic, err
	// // }

	return topic, err

}

func PostTopics (w http.ResponseWriter, r *http.Request, user_id int) (err error) {
	
	err = godotenv.Load(".env")
	if err != nil {
		fmt.Printf("読み込み出来ませんでした: %v", err)
	} 
	Db, err := gorm.Open(postgres.Open(os.Getenv("DB_DSN")), &gorm.Config{})
	if err != nil {
	panic("failed to connect database")
	}

	if r.Method == "POST" {
		topic := BbsTopic{
  	  TopicTitle: r.FormValue("title"),
  	  TopicDescription: r.FormValue("description"),
			TopicCategory: r.FormValue("category"),
			BbsUserId: user_id,
		}
	 
		topic_pointer := Db.Create(&topic)
		if topic_pointer.Error != nil {
  	  panic("failed to insert record")
		}
	}

	return err
}

func GetTopics (w http.ResponseWriter, r *http.Request) (Topics []BbsTopic, err error) {
	
	// TOOD:この処理をまとめて関数化できないか検討
	Db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer Db.Close()

	// TODO:Goで動的SQLを実装する仕組みがあればそちらで実装したい
	title := r.FormValue("title")
	description := r.FormValue("description")
	category := r.FormValue("category")

	var conditions []string
	var params []interface{}

	if title != "" {
		conditions = append(conditions, "topic_title ILIKE $"+strconv.Itoa(len(conditions)+1))
		params = append(params, "%"+title+"%")
	}
	if description != "" {
		conditions = append(conditions, "topic_description ILIKE $"+strconv.Itoa(len(conditions)+1))
		params = append(params, "%"+description+"%")
	}
	if category != "" {
		conditions = append(conditions, "topic_category ILIKE $"+strconv.Itoa(len(conditions)+1))
		params = append(params, "%"+category+"%")
	}

	cmd := "SELECT id, topic_title, topic_description, topic_category, created_at FROM bbs_topics WHERE " + strings.Join(conditions, " OR ")
	
	rows, err := Db.Query(cmd, params...)
	if err != nil {
		log.Fatalln(err)
	}
	for rows.Next() {
		var topic BbsTopic

		// weekday := strings.NewReplacer(
		// 	"Sun", "日",
		// 	"Mon", "月",
		// 	"Tue", "火",
		// 	"Wed", "水",
		// 	"Thu", "木",
		// 	"Fri", "金",
		// 	"Sat", "土",
		// )

		err = rows.Scan(
			&topic.ID,
			&topic.TopicTitle,
			&topic.TopicDescription,
			&topic.TopicCategory,
			&topic.CreatedAt,
		)
		if err != nil {
			log.Fatalln(err)
		}
		
		// formattedDate := weekday.Replace(createdAt.Format("2006年1月2日(Mon) 15:04:05"))
		// topic.CreatedAt = formattedDate

		Topics = append(Topics, topic)
	}

	rows.Close()

	return Topics, err
}
