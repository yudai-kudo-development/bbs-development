package models

import (
	//"database/sql"
	"fmt"
	//"log"
	"net/http"
	// "strconv"
	// "strings"
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
  	BbsUserId           int
	CreatedAt        time.Time
	UpdatedAt        time.Time
	BbsReplies       []BbsReply
	User             BbsUser `gorm:"foreignKey:BbsUserId"`
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

func GetIndividualTopic (id int) (BbsTopic, error) {

	err = godotenv.Load(".env")
	if err != nil {
		fmt.Printf("読み込み出来ませんでした: %v", err)
	} 
	Db, err := gorm.Open(postgres.Open(os.Getenv("DB_DSN")), &gorm.Config{})
	if err != nil {
	panic("failed to connect database")
	}

	var topic BbsTopic
	err = Db.Preload("BbsReplies").First(&topic, id).Error
	fmt.Println(topic)

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
	
	err = godotenv.Load(".env")
	if err != nil {
		fmt.Printf("読み込み出来ませんでした: %v", err)
	} 
	Db, err := gorm.Open(postgres.Open(os.Getenv("DB_DSN")), &gorm.Config{})
	if err != nil {
	panic("failed to connect database")
	}

	// TODO:Goで動的SQLを実装する仕組みがあればそちらで実装したい
	title := r.FormValue("title")
	description := r.FormValue("description")
	category := r.FormValue("category")

	var topics []BbsTopic
	query := Db.Where("1=1")
	fmt.Printf("%v", query)

	if title != "" {
		query = query.Where("topic_title LIKE ?", "%"+title+"%")
	}
	if description != "" {
		query = query.Where("topic_description LIKE ?", "%"+description+"%")
	}
	if category != "" {
		query = query.Where("topic_category LIKE ?", "%"+category+"%")
	}

	query.Find(&topics)

	return topics, err
}

func UpdateTopic (w http.ResponseWriter, r *http.Request, topic_id int) (err error) {
	
	err = godotenv.Load(".env")
	if err != nil {
		fmt.Printf("読み込み出来ませんでした: %v", err)
	} 
	Db, err := gorm.Open(postgres.Open(os.Getenv("DB_DSN")), &gorm.Config{})
	if err != nil {
	panic("failed to connect database")
	}

	// TODO:Goで動的SQLを実装する仕組みがあればそちらで実装したい
	topicID := r.FormValue("id")
	title := r.FormValue("title")
	description := r.FormValue("description")
	category := r.FormValue("category")

	// 更新内容を保持するマップを準備
	updateData := make(map[string]interface{})
	if title != "" {
		updateData["topic_title"] = title
	}
	if description != "" {
		updateData["topic_description"] = description
	}
	if category != "" {
		updateData["topic_category"] = category
	}

	// 更新処理
	if len(updateData) > 0 {
		err = Db.Model(&BbsTopic{}).Where("id = ?", topicID).Updates(updateData).Error
		if err != nil {
			fmt.Printf("更新エラー: %v", err)
		}
	}

	return err
}

func DeleteTopic (w http.ResponseWriter, r *http.Request, topic_id int) (err error) {
	
	err = godotenv.Load(".env")
	if err != nil {
		fmt.Printf("読み込み出来ませんでした: %v", err)
	} 
	Db, err := gorm.Open(postgres.Open(os.Getenv("DB_DSN")), &gorm.Config{})
	if err != nil {
	panic("failed to connect database")
	}

	var topics BbsTopic
	Db.Where("id = ?", topic_id).Delete(&topics)

	return err
}
