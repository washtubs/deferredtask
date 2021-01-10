package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"

	"github.com/washtubs/assign2me"
)

func main() {
	svc := assign2me.GetDeferrableService()
	flag.Parse()
	action := flag.Arg(0)
	if action == "" {
		log.Fatal("Must provide an action: do | add | ls | dismiss ")
	}
	if action == "do" {
		usage := "assign2me do <idx>"
		idStr := flag.Arg(1)
		if idStr == "" {
			log.Fatal(usage)
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Fatal(err)
		}

		err = svc.DoTask(id)
		if err != nil {
			log.Fatal(err)
		}

	} else if action == "add" {
		usage := "assign2me add <description> <cmd>"
		description := flag.Arg(1)
		cmd := flag.Arg(2)
		if cmd == "" || description == "" {
			log.Fatal(usage)
		}
		err := svc.AddTask(assign2me.DeferrableTask{
			Description: description,
			Cmd:         cmd,
		})
		if err != nil {
			log.Fatal(err)
		}

	} else if action == "ls" {
		tasks, err := svc.ListTasks()
		if err != nil {
			log.Fatal(err)
		}
		for idx, task := range tasks {
			fmt.Printf("%d\t%s\t%s\n", idx, task.Description, task.Cmd)
		}
	} else if action == "dismiss" {
		usage := "assign2me dismiss <idx>"
		idStr := flag.Arg(1)
		if idStr == "" {
			log.Fatal(usage)
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Fatal(err)
		}

		err = svc.DismissTask(id)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("Must provide an action: do | add | ls | dismiss ")
	}
}
