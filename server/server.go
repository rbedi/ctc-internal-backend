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

type Project struct {
	Id int
	Title string
	Github string
	Organization string
	Description string
}

var orm beedb.Model

func main() {
	db, err := sql.Open("sqlite3", "./main.db")
	if err != nil {
		panic(err)
	}
	orm = beedb.New(db)
	fmt.Printf("Here are the tags with their ids:\n")
	printTags()
	var projname string
	fmt.Printf("Add a new project. Name: ")
	fmt.Scanf("%s", &projname)
	fmt.Printf("Associate a tag with your project by id: ")
	var projtag int
	fmt.Scanf("%d", %projtag)

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


/* func getTag() {

	var mytag Tags


} */
