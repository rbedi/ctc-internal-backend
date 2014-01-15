package main

import (
	"fmt"
	"github.com/astaxie/beedb"
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
)

type Tag struct {
	Uid int
	Name string
}

var orm beedb.Model

func main() {
	db, err := sql.Open("sqlite3", "./main.db")
	if err != nil {
		panic(err)
	}
	orm = beedb.New(db)
	insertTag()
}

func insertTag() {
	var mytag Tag
	mytag.Name = "First Tag"
	orm.Save(&mytag)
	fmt.Println(mytag)
}
