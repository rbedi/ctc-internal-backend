package main

import (
	"fmt"
	"github.com/astaxie/beedb"
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
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
	Title string
	Github string
	Organization string
	Description string
}

type ProjectTag struct {
	ProjectId int
	TagId int
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
	fmt.Printf("Would you like to make a new project? (true or false)")
	var makeProj bool
	fmt.Scanf("%t", &makeProj)
	if (makeProj) {
		addProject()
	}
}

func addProject() {
	var newProject Project
	fmt.Printf("Add a new project. Title: ")
	fmt.Scanf("%s", &newProject.Title)
	fmt.Printf("Github link: ")
	fmt.Scanf("%s", &newProject.Github)
	fmt.Printf("Organization: ")
	fmt.Scanf("%s", &newProject.Organization)
	fmt.Printf("Description: ")
	fmt.Scanf("%s", &newProject.Description)

	insertProject(newProject)

	fmt.Printf("Would you like to associate a tag with your project? (true or false)")
	var doAssociate bool
	fmt.Scanf("%t", &doAssociate)
	if (doAssociate) {
		associateTag()
	}
}

func associateTag() {
	fmt.Printf("Associate a tag with your project by tag id: ")
	var projtag ProjectTag
	fmt.Scanf("%d", &projtag.TagId)
}

func insertProject(newProject Project) error {
	return orm.Save(&newProject)
}

func insertTag() {
	var mytag Tag
	mytag.Name = "First Tag"
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

func printTags() {
	var allTags []Tag = getAllTags()
	fmt.Println(allTags)
} 

func getAllTags() []Tag {
	var allTags []Tag
	orm.FindAll(&allTags)
	return allTags
}

func getProjectsWithTag(tagId int) []Project {
	var projectsWithTag []Project
	orm.Where("projectId = ?",tagId).FindAll(&projectsWithTag)
	return projectsWithTag
}

