package main

import (
	"fmt"
)

func main() {
	initDB()
	fmt.Printf("Here are the tags with their ids:\n")
	printTags()
	fmt.Printf("Would you like to make a new project? (true or false) ")
	var makeProj bool
	fmt.Scanf("%t\n", &makeProj)
	if (makeProj) {
		addProject()
	}
}