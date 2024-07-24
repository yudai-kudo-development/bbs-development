package models

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type User struct {
	ID        int
	UUID      string
	Name      string
	Email     string
	PassWord  string
	CreatedAt time.Time
	Topics    []Topic
}

type Session struct {
	ID        int
	UUID      string
	Email     string
	UserID    int
	CreatedAt time.Time
}

func (u *User) CreateUser() (err error) {

	Db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer Db.Close()

	cmd := `insert into bbs_users (
	uuid,
	name,
	email,
	password,
	created_at ) values ($1, $2, $3, $4, $5)`

	_, err = Db.Exec(cmd, createUUID(), u.Name, u.Email, Encrypt(u.PassWord),time.Now())

	if err != nil {
		log.Fatalln(err)
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
	
	fmt.Println("DB挿入処理完了")
	fmt.Println(err)
	fmt.Println(sess.ID)

	if err != nil {
		valid = false
		return
	}

	if sess.ID != 0 {
		valid = true
	}

	return valid, err
}

func GetUserByEmail (email string) (user User,err error) {
	Db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer Db.Close()

	user = User{}
	cmd := `select id, uuid, name, email, password, created_at
	from bbs_users where email = $1`
	err = Db.QueryRow(cmd, email).Scan(&user.ID, &user.UUID,&user.Name,&user.Email,&user.PassWord,&user.CreatedAt)

	return user,err
}

func (u *User) CreateSession () (session Session,err error) {
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

func (sess *Session) GetUserBySession() (user User, err error) {

	Db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer Db.Close()

	user = User{}
	cmd := `select id, uuid, name, email, created_at
	from bbs_users where id = $1`
	err = Db.QueryRow(cmd, sess.UserID).Scan(&user.ID, &user.UUID, &user.Name, &user.Email, &user.CreatedAt)
	
	return user,err
}

func GetUserWithTopics (user User, id int) (User, error) {

	Db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer Db.Close()

	bbs_topics := `SELECT id, topic_title, topic_description, topic_category, created_at FROM bbs_topics WHERE id = $1`
	rows, err := Db.Query(bbs_topics, id)
	if err != nil {
			return user, err
	}
	defer rows.Close()

	for rows.Next() {
			var topic Topic
			err := rows.Scan(&topic.ID, &topic.Title, &topic.Description, &topic.Category, &topic.CreatedAt)
			if err != nil {
					return user, err
			}
			user.Topics = append(user.Topics, topic)
	}

	if err := rows.Err(); err != nil {
		return user, err
	}

	return user, err
}
