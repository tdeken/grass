package main

import (
	"flag"
	"fmt"
	"github.com/tdeken/grass/createproto"
	"github.com/tdeken/grass/fiberaction"
	"github.com/tdeken/grass/fiberweb"
	"os"
)

type Build interface {
	Run()
	Error() error
}

const (
	sceneCreateNone = iota
	sceneFiberWebInit
	sceneFiberWeb
	sceneProto
)

type Args struct {
	scene        int    // what you  will do scene
	FiberWebInit string // fiber web frame init
	FiberWeb     string // build fiber action
	Proto        string // build proto
	Dir          string // root dir
}

var help bool
var params = &Args{}

var fs = flag.NewFlagSet("grass", flag.ExitOnError)

func init() {
	defer parse()

	// help
	flag.BoolVar(&help, "h", false, "grass options help")
	flag.BoolVar(&help, "help", false, "grass options help")

	// create fiber web frame init
	fs.StringVar(&params.FiberWebInit, "fbinit", "", "create fiber web frame")

	// build fiber action
	fs.StringVar(&params.FiberWeb, "fb", "", "build fiber action \r\n es: -fb [dir]")

	// build fiber action
	fs.StringVar(&params.Proto, "bp", "", "build proto name")
	fs.StringVar(&params.Dir, "d", "", "root dir name")
}

// 解析命令行参数
func parse() {
	var args = os.Args[1:]
	if err := fs.Parse(args); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if params.FiberWebInit != "" {
		params.scene = sceneFiberWebInit
	} else if params.FiberWeb != "" {
		params.scene = sceneFiberWeb
	} else if params.Proto != "" {
		params.scene = sceneProto
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
	case sceneFiberWebInit:
		build = fiberweb.NewFiberWeb(params.FiberWebInit)
	case sceneFiberWeb:
		build = fiberaction.NewFiberAction(params.FiberWeb)
	case sceneProto:
		build = createproto.NewCreateProto(params.Dir, params.Proto)
	}

	if build == nil {
		return
	}

	build.Run()
	if build.Error() != nil {
		fmt.Printf("%v \r\n", build.Error())
		return
	}
}
