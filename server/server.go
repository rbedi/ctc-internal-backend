package main

import (
	"fmt"
	"github.com/astaxie/beedb"
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
)

type Tag struct {
	Id int
	Title string
}

var orm beedb.Model

func main() {
	db, err := sql.Open("sqlite3", "./main.db")
	if err != nil {
		panic(err)
	}

	beedb.PluralizeTableNames=true
	orm = beedb.New(db)
	insertTag()
	getAllTags()
}

func insertTag() {
	var mytag Tag
	mytag.Title = "First Tag"
	err := orm.Save(&mytag)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(mytag)
	}
}


func getTag(tagId int) {
	var mytag Tag
	orm.Where("id=?",tagId).Find(&mytag)
	fmt.Println(mytag)
} 

func getAllTags(){
	var allTags []Tag
	orm.FindAll(&allTags)
	fmt.Println(allTags)
}

func getProjectsWithTag(tagId int){
	var projectsWithTag Project
	orm.Where("projectId = ?",tagId).FindAll(&projectsWithTag)
	fmt.Println(projectsWithTag)
}

