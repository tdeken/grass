package main

import (
	"flag"
	"fmt"
	"grass/createproject"
	"os"
)

const (
	sceneCreateProject = iota + 1
)

type Args struct {
	scene         int    // what you  will do scene
	CreateProject string // project_name to be created
}

var params = &Args{}

var fs = flag.NewFlagSet("grass", flag.ExitOnError)

func init() {
	defer parse()

	// create_project
	fs.StringVar(&params.CreateProject, "cp", "", "project name to be created")
	fs.StringVar(&params.CreateProject, "create_project", "", "project name to be created")
}

func help() {

}

// 解析命令行参数
func parse() {
	if err := fs.Parse(os.Args[1:]); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if params.CreateProject != "" {
		params.scene = sceneCreateProject
	}
}

func main() {

	switch params.scene {
	case sceneCreateProject:
		cp := createproject.NewCreateProject(params.CreateProject)

		cp.Run()
		if cp.Error() != nil {
			fmt.Println(cp.Error())
			os.Exit(1)
		}
	}
}
