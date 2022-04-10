package database

import (
	"github.com/gogf/gf/database/gdb"
	"github.com/skiy/gfutils/ldb"
)

// GetUcDB GetUcDB
func GetUcDB() (db gdb.DB) {
	db = ldb.Get("uc")
	return db
}

// GetXiunoDB GetXiunoDB
func GetXiunoDB() (db gdb.DB) {
	db = ldb.Get("xiuno")
	return db
}

// GetDiscuzDB GetDiscuzDB
func GetDiscuzDB() (db gdb.DB) {
	db = ldb.Get("discuz")
	return db
}
