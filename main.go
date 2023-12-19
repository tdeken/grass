package main

import (
	"flag"
	"fmt"
	"github.com/tdeken/grass/createproject"
	"os"
)

const (
	sceneCreateNone = iota
	sceneCreateProject
)

type Args struct {
	scene         int    // what you  will do scene
	CreateProject string // project_name to be created
}

var help bool
var params = &Args{}

var fs = flag.NewFlagSet("grass", flag.ExitOnError)

func init() {
	defer parse()

	// help
	flag.BoolVar(&help, "h", false, "grass options help")
	flag.BoolVar(&help, "help", false, "grass options help")

	// create_project
	fs.StringVar(&params.CreateProject, "cp", "", "project name to be created")
	fs.StringVar(&params.CreateProject, "create_project", "", "project name to be created")
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

	if help || params.scene == sceneCreateNone {
		_, _ = fmt.Fprintln(os.Stderr, "grass usage options:")
		fs.PrintDefaults()
		os.Exit(0)
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
