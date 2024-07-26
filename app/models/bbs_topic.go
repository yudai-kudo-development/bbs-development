package models

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

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

func GetIndividualTopic (id int) (Topic, error) {
	
	var topic Topic

	Db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer Db.Close()

	bbs_topic := `SELECT id, topic_title, topic_description, topic_category, created_at FROM bbs_topics WHERE id = $1`
	err = Db.QueryRow(bbs_topic, id).Scan(&topic.ID, &topic.Title, &topic.Description, &topic.Category, &topic.CreatedAt)
	if err != nil {
			return topic, err
	}

	bbs_replies := `SELECT id, bbs_topic_id, reply_name, reply_content, created_at FROM bbs_replies WHERE bbs_topic_id = $1`
	rows, err := Db.Query(bbs_replies, id)
	if err != nil {
			return topic, err
	}
	defer rows.Close()

	for rows.Next() {
			var reply Reply
			err := rows.Scan(&reply.ID, &reply.TopicID, &reply.Name, &reply.Content, &reply.CreatedAt)
			if err != nil {
					return topic, err
			}
			topic.Replies = append(topic.Replies, reply)
	}

	if err := rows.Err(); err != nil {
		return topic, err
	}

	return topic, err
}

func PostTopics (w http.ResponseWriter, r *http.Request, user_id *int ) (err error) {
	Db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer Db.Close()
	
	if r.Method == "POST" {
		userId := user_id
		title := r.FormValue("title")
		description := r.FormValue("description")
		category := r.FormValue("category")
	
		insert, err := Db.Prepare("INSERT INTO bbs_topics(topic_title, topic_description, topic_category, bbs_user_id) VALUES ($1,$2,$3,$4)")
		if err != nil {
			fmt.Println(err)
		}
	
		insert.Exec(title, description, category, userId)
	} 

	return err
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
