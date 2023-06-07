package main

import (
	"Mod30/pkg/storage"
	"Mod30/pkg/storage/postgres"
	"fmt"
	"log"
)

var db storage.DBInterface

func main() {
	var err error
	var task postgres.Task
	pwd := "******"
	connstr := "postgres://postgres:" + pwd + "@192.168.5.136/postgres"
	db, err = postgres.New(connstr)
	if err != nil {
		log.Fatal(err)
	}
	//task.AssignedID = 2
	//task.AuthorID = 1
	task.ID = 3
	task.Content = "разработать 1"
	task.Title = "задача 1"
	err = db.UpdateTask(task)
	//inti, err := db.NewTask(task)
	fmt.Println(err)
	tasks, err := db.Tasks(0, 0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(tasks)

}
