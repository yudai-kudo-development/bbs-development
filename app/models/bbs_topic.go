package models

import (
	"fmt"
	// "log"
	"net/http"
	"strings"
	// "todo_app/app/models"
	// "path/filepath"
	"database/sql"
	"log"
	"time"
	"strconv"

	_ "github.com/lib/pq"
)

type Topic struct {
    ID          int
    Title       string
    Description string
	Category   string
    CreatedAt   string
    Replies     []Reply
}

func GetRecentTopics (w http.ResponseWriter, r *http.Request) (Topics []Topic, err error) {
	
	// TOOD:この処理をまとめて関数化できないか検討
	Db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer Db.Close()

	cmd := `SELECT id, topic_title, topic_description, topic_category, created_at FROM bbs_topics ORDER BY created_at DESC LIMIT 5`
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
			&topic.ID,
			&topic.Title,
			&topic.Description,
			&topic.Category,
			&createdAt,
		)
		if err != nil {
			log.Fatalln(err)
		}
		
		formattedDate := weekday.Replace(createdAt.Format("2006年1月2日(Mon) 15:04:05"))
		topic.CreatedAt = formattedDate

		Topics = append(Topics, topic)
	}

	rows.Close()

	return Topics, err
}

func GetIndividualTopic (id int) ([]Topic, error) {
	
	topics := []Topic{}
    topicMap := make(map[int]*Topic)
	// var topicCreatedAt time.Time
	// var replyCreatedAt time.Time

	Db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer Db.Close()

	// weekday := strings.NewReplacer(
	// 	"Sun", "日",
	// 	"Mon", "月",
	// 	"Tue", "火",
	// 	"Wed", "水",
	// 	"Thu", "木",
	// 	"Fri", "金",
	// 	"Sat", "土",
	// )

	cmd := `SELECT * FROM topic_view WHERE topic_id = $1`
	rows, err := Db.Query(cmd, id)
	if err != nil {
		log.Fatalln(err)
	}

	for rows.Next() {
		var (
            topicID         int
            topicTitle      string
            topicDescription string
			topicCategory   string
            topicCreatedAt  string
            replyID         sql.NullInt64
			replyTopicID    sql.NullInt64
            replyName       sql.NullString
            replyContent    sql.NullString
            replyCreatedAt  sql.NullString
        )

        err := rows.Scan(&topicID, &topicTitle, &topicDescription, &topicCategory, &topicCreatedAt, &replyID, &replyTopicID, &replyName, &replyContent, &replyCreatedAt)
        if err != nil {
            return nil, err
        }

        if _, exists := topicMap[topicID]; !exists {
            topic := Topic{
                ID:          topicID,
                Title:       topicTitle,
                Description: topicDescription,
				Category:    topicCategory,
                CreatedAt:   topicCreatedAt,
                Replies:     []Reply{},
            }
            topics = append(topics, topic)
            topicMap[topicID] = &topics[len(topics)-1]
        }

		if replyID.Valid {
            reply := Reply{
                ID:        int(replyID.Int64),
                TopicID:   topicID,
                Name:      replyName.String,
                Content:   replyContent.String,
                CreatedAt: replyCreatedAt.String,
            }
            topicMap[topicID].Replies = append(topicMap[topicID].Replies, reply)
        }
		
		// formattedTopicDate := weekday.Replace(topicCreatedAt.Format("2006/1/2(Mon) 15:04:05"))
		// IndividualTopic.TopicCreatedAt = formattedTopicDate
		// formattedReplyDate := weekday.Replace(replyCreatedAt.Format("2006/1/2(Mon) 15:04:05"))
		// IndividualTopic.ReplyCreatedAt = formattedReplyDate
	}
	
	defer rows.Close()

	return  topics, nil
}

func PostTopics (w http.ResponseWriter, r *http.Request) (err1 error, err2 error) {
	Db, err1 := sql.Open("postgres", connStr)
	if err1 != nil {
		log.Fatal(err1)
	}
	defer Db.Close()
	
	if r.Method == "POST" {
		title := r.FormValue("title")
		description := r.FormValue("description")
		category := r.FormValue("category")
	
		insert, err2 := Db.Prepare("INSERT INTO bbs_topics(topic_title, topic_description, topic_category) VALUES ($1,$2,$3)")
		if err2 != nil {
			fmt.Println(err2)
		}
	
		insert.Exec(title, description, category)
	} 

	return err1,err2
}

func GetTopics (w http.ResponseWriter, r *http.Request) (Topics []Topic, err error) {
	
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
			&topic.ID,
			&topic.Title,
			&topic.Description,
			&topic.Category,
			&createdAt,
		)
		if err != nil {
			log.Fatalln(err)
		}
		
		formattedDate := weekday.Replace(createdAt.Format("2006年1月2日(Mon) 15:04:05"))
		topic.CreatedAt = formattedDate

		Topics = append(Topics, topic)
	}

	rows.Close()

	return Topics, err
}
