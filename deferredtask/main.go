package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/washtubs/deferredtask"
)

func main() {
	svc := deferredtask.GetDeferrableService()
	fs := flag.NewFlagSet("deferred-task", flag.ExitOnError)
	fs.Parse(os.Args[1:])
	//flag.Parse()
	action := fs.Arg(0)
	if action == "" {
		log.Fatal("Must provide an action: do | add | ls | dismiss ")
	}
	actionArgs := fs.Args()[1:]
	log.Printf("%+v", actionArgs)
	if action == "do" {
		usage := "deferred-task do <idx>"
		fs = flag.NewFlagSet("deferred-task", flag.ExitOnError)
		fs.Parse(actionArgs)
		idStr := fs.Arg(0)
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
		fs = flag.NewFlagSet("deferred-task", flag.ExitOnError)
		handle := ""
		fs.StringVar(&handle, "handle", "", "Handle to use to update an existing value if it exists")
		err := fs.Parse(actionArgs)
		if err != nil {
			log.Fatal(err)
		}

		usage := "deferred-task add <options> <description> <cmd>"
		description := fs.Arg(0)
		cmd := fs.Arg(1)
		if cmd == "" || description == "" {
			log.Fatal(usage)
		}
		err = svc.AddTask(deferredtask.DeferrableTask{
			Description: description,
			Cmd:         cmd,
			Handle:      handle,
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
			fmt.Printf("%d\t%s\t%s\t%s\n", idx, task.Description, task.Cmd, task.Handle)
		}
	} else if action == "dismiss" {
		fs = flag.NewFlagSet("deferred-task", flag.ExitOnError)
		fs.Parse(actionArgs)
		usage := "deferred-task dismiss <idx>"
		idStr := fs.Arg(0)
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
