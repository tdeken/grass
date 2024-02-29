package main

import (
	"flag"
	"fmt"
	"github.com/tdeken/grass/fiberweb"
	"os"
)

type Build interface {
	Run()
	Error() error
}

const (
	sceneCreateNone = iota
	sceneFiberWeb
)

type Args struct {
	scene    int    // what you  will do scene
	FiberWeb string // fiber web frame
}

var help bool
var params = &Args{}

var fs = flag.NewFlagSet("grass", flag.ExitOnError)

func init() {
	defer parse()

	// help
	flag.BoolVar(&help, "h", false, "grass options help")
	flag.BoolVar(&help, "help", false, "grass options help")

	// create fiber web frame
	fs.StringVar(&params.FiberWeb, "fb", "", "create fiber web frame")
}

// 解析命令行参数
func parse() {
	var args = os.Args[1:]
	if err := fs.Parse(args); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if params.FiberWeb != "" {
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
	case sceneFiberWeb:
		build = fiberweb.NewFiberWeb(params.FiberWeb)
	}

	if build == nil {
		return
	}

	build.Run()
	if build.Error() != nil {
		fmt.Printf("%v \r", build.Error())
		return
	}
}
