package main

import (
	"fmt"
	"net/http"
	"strconv"
	"encoding/json"
	_"bytes"
)


func projectInfoHandler (rw http.ResponseWriter, req *http.Request){
    rw.Header().Set("Content-Type", "application/json")
	projectId, err := strconv.Atoi(req.URL.Path[9:])
	if err != nil{
		fmt.Println(err)
	}
	proj := getProjectInfo(projectId)
	b, err := json.Marshal(proj)
	if err != nil {
		fmt.Println(err)
	}
    s := string(b[:])
	fmt.Fprintf(rw,s)
}

/*func addProjectHandler (rw http.ResponseWriter, req *http.Request){


}*/

func main() {

	initDB()

	http.HandleFunc("/project/", projectInfoHandler)
	//http.HandleFunc("/project", addProjectHandler)
	http.ListenAndServe(":8080", nil)

	fmt.Printf("Here are the tags with their ids:\n")
	printTags()
	fmt.Printf("Would you like to make a new project? (true or false) ")
	var makeProj bool
	fmt.Scanf("%t\n", &makeProj)
	if (makeProj) {
		addProject()
	}


}

