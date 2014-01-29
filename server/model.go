package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	_ "database/sql"
)

type Tag struct {
	Id int
	Name string
}

func (tag Tag) String() string {
	return fmt.Sprintf("\n\t{\n\t id: %d\n\t name: %s\n\t}\n", tag.Id, tag.Name)
}

type Project struct {
	Id int
	Title string `db:"title"`
	Github string `db:"github"`
	Organization string `db:"organization"`
	Description string `db:"description"`
}

type ProjectTag struct {
	ProjectId int
	TagId int
}

var db *sqlx.DB

func initDB() {
	var err error
	db, err = sqlx.Open("sqlite3", "./main.db")
	if err != nil {
		panic(err)
	}
}

func addProject() {
	var newProject Project
	fmt.Printf("Add a new project Title: ")
	fmt.Scanf("%s\n", &newProject.Title)
	fmt.Printf("Github link: ")
	fmt.Scanf("%s\n", &newProject.Github)
	fmt.Print("Organization: ")
	fmt.Scanf("%s\n", &newProject.Organization)
	fmt.Print("Description: ")
	fmt.Scanf("%s\n", &newProject.Description)

	insertProject(newProject)

	fmt.Printf("Would you like to associate a tag with your project? (true or false) ")
	var doAssociate bool
	fmt.Scanf("%t\n", &doAssociate)
	if (doAssociate) {
		associateTag()
	}
}

func associateTag() {
	fmt.Printf("Associate a tag with your project by tag id: ")
	var projtag ProjectTag
	fmt.Scanf("%d\n", &projtag.TagId)

}

func insertProject(newProject Project) error {
	tx := db.MustBegin()
	tx.NamedExec("INSERT INTO project (title, github, organization, description) VALUES (:title, :github, :organization, :description)", &newProject)
	err := tx.Commit()
	return err
}

func printTags() {
	var allTags []Tag = getAllTags()
	fmt.Println(allTags)
} 

func getAllTags() []Tag {
	var allTags []Tag
	db.Select(&allTags, "SELECT * FROM tag")
	return allTags
}