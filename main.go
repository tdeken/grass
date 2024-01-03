package main

import (
	"flag"
	"fmt"
	"github.com/tdeken/grass/createproject"
	"os"
)

type Build interface {
	Run()
	Error() error
}

const (
	sceneCreateNone = iota
	sceneCreateProject
	sceneFiberWeb
)

type Args struct {
	scene         int    // what you  will do scene
	CreateProject string // project_name to be created
	FiberWeb      string // fiber web frame
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

	// create fiber web frame
	fs.StringVar(&params.CreateProject, "fb", "", "project name to be created")
}

// 解析命令行参数
func parse() {
	var args = os.Args[1:]
	if err := fs.Parse(args); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if params.CreateProject != "" {
		params.scene = sceneCreateProject
	} else if params.FiberWeb != "" {
		params.scene = sceneFiberWeb
	}

	if help || params.scene == sceneCreateNone {
		_, _ = fmt.Fprintln(os.Stderr, "grass usage options:")
		fs.PrintDefaults()
		os.Exit(0)
	}
}

func main() {
	var build Build
	switch params.scene {
	case sceneCreateProject:
		build = createproject.NewCreateProject(params.CreateProject)
		build.Run()
	case sceneFiberWeb:

	}

	if build != nil && build.Error() != nil {
		fmt.Printf("%v \r\n", build.Error())
		os.Exit(1)
	}
}
