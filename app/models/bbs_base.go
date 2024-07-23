package models

import (
	"fmt"
	"github.com/google/uuid"
	"crypto/sha1"
	"database/sql"
)

var err error
var Db *sql.DB
// TODO:configに接続情報をまとめる
var connStr = "user=yudai.kudo dbname=bbs_development sslmode=disable"

func createUUID() (uuidobj uuid.UUID) {
	uuidobj, _ = uuid.NewUUID()
	return uuidobj
}

func Encrypt(plaintext string) (cryptext string) {
	cryptext = fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
	return cryptext
}
