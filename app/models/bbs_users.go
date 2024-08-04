package models

import (
	"database/sql"
	"log"
	"time"
	"fmt"
	"os"
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type BbsUser struct {
	ID        int
	UUID      string
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
	BbsTopics []BbsTopic `gorm:"foreignKey:BbsUserId"`
}

type Session struct {
	ID        int
	UUID      string
	Email     string
	UserID    int
	CreatedAt time.Time
}

func (u *BbsUser) CreateUser() (err error) {

	err = godotenv.Load(".env")
	if err != nil {
		fmt.Printf("読み込み出来ませんでした: %v", err)
	} 
	Db, err := gorm.Open(postgres.Open(os.Getenv("DB_DSN")), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	user := BbsUser{
		Name:      u.Name,
		Email:     u.Email,
		Password:  Encrypt(u.Password),
		CreatedAt: time.Now(),
	}

	topic_pointer := Db.Create(&user)
	if topic_pointer.Error != nil {
		panic("failed to insert record")
	}

	return err
}

func (sess *Session) CheckSession() (valid bool, err error) {

	Db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer Db.Close()

	cmd := `select id, uuid, email, user_id, created_at
	from bbs_sessions where uuid = $1`

	err = Db.QueryRow(cmd, sess.UUID).Scan(
		&sess.ID, &sess.UUID,&sess.Email, &sess.UserID, &sess.CreatedAt)

	if err != nil {
		valid = false
		return
	}

	if sess.ID != 0 {
		valid = true
	}

	return valid, err
}

func GetUserByEmail (email string) (user BbsUser,err error) {
	Db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer Db.Close()

	user = BbsUser{}
	cmd := `select id, uuid, name, email, password, created_at
	from bbs_users where email = $1`
	err = Db.QueryRow(cmd, email).Scan(&user.ID, &user.UUID,&user.Name,&user.Email,&user.Password,&user.CreatedAt)

	return user,err
}

func (u *BbsUser) CreateSession () (session Session,err error) {
	Db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer Db.Close()

	session = Session{}
	cmd1 := `insert into bbs_sessions (
	uuid,
	email,
	user_id,
	created_at) values ($1, $2, $3, $4)`

	_, err = Db.Exec(cmd1, createUUID(), u.Email, u.ID, time.Now())
	if err != nil {
		log.Fatalln(err)
	}

	cmd2 := `select id, uuid, email, user_id, created_at
	from bbs_sessions where user_id = $1 and email = $2`

	err = Db.QueryRow(cmd2, u.ID , u.Email).Scan(&session.ID, &session.UUID, &session.Email, &session.UserID,&session.CreatedAt)

	return session,err
}

func (sess *Session) GetUserBySession() (user BbsUser, err error) {

	Db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer Db.Close()

	user = BbsUser{}
	cmd := `select id, uuid, name, email, created_at
	from bbs_users where id = $1`
	err = Db.QueryRow(cmd, sess.UserID).Scan(&user.ID, &user.UUID, &user.Name, &user.Email, &user.CreatedAt)
	
	return user,err
}

func GetTopicsWithUser (bbs_user_id int) (users []BbsUser, err error) {

	err = godotenv.Load(".env")
	if err != nil {
		fmt.Printf("読み込み出来ませんでした: %v", err)
	} 
	Db, err := gorm.Open(postgres.Open(os.Getenv("DB_DSN")), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

    err = Db.Model(&BbsUser{}).Preload("BbsTopics").Where("id = ?", bbs_user_id).Find(&users).Error

	return users, err
}
