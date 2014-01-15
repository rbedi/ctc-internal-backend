package main

import (
	"fmt"
	"github.com/astaxie/beedb"
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
)

type Tags struct {
	Id int
	Title string
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
	var mytag Tags
	mytag.Title = "First Tag"
	err := orm.Save(&mytag)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(mytag)
	}
}
