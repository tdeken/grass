package fiberaction

import (
	"errors"
	"fmt"
	"github.com/tdeken/grass/basic"
	"github.com/tdeken/grass/utils"
	"gopkg.in/yaml.v3"
	"os"
)

type FiberAction struct {
	basic.Basic
	conf basic.GrassConf
}

func (s *FiberAction) Error() (err error) {
	if s.Err == nil {
		return
	}

	return errors.New(fmt.Sprintf("FiberAction err: %v", s.Err))
}

// NewFiberAction create_project
func NewFiberAction(modName string) *FiberAction {
	fa := new(FiberAction)
	{
		fa.Init(modName)
	}

	return fa
}

func (s *FiberAction) Run() {
	var err error
	defer func() {
		s.Err = err
	}()

	err = s.loadConf()
	if err != nil {
		return
	}

	err = s.dir()
	if err != nil {
		return
	}

	return
}

func (s *FiberAction) loadConf() (err error) {
	cfd := s.PrefixDir(fmt.Sprintf("etc/grass.yaml"))

	in, err := os.ReadFile(cfd)
	if err != nil {
		err = errors.New(fmt.Sprintf("读取grass.yaml失败: %v", err))
		return
	}

	err = yaml.Unmarshal(in, &s.conf)
	if err != nil {
		err = errors.New(fmt.Sprintf("读取grass.yaml参数异常: %v", err))
		return
	}

	return
}

func (s *FiberAction) dir() (err error) {
	serDir, hdlDir, srsDir := s.PrefixDir(s.conf.Analyze.Service), s.PrefixDir(s.conf.Analyze.Handler), s.PrefixDir(s.conf.Analyze.Sources)

	err = utils.MkDirAll(serDir, hdlDir, srsDir)
	if err != nil {
		return
	}

	serPkg, serStruct := utils.PkgAndStruct(s.conf.Analyze.Service)
	ser, err := utils.CreateTmp(ServiceTemp{
		Pkg:     serPkg,
		Service: serStruct,
	}, serviceTemp)
	if err != nil {
		return
	}

	serFile := fmt.Sprintf("%s/%s.go", serDir, serPkg)
	err = utils.NotExistCreateFile(serFile, ser)
	if err != nil {
		return
	}

	hdlPkg, hdlStruct := utils.PkgAndStruct(s.conf.Analyze.Handler)
	hdl, err := utils.CreateTmp(HandlerTemp{
		ModName: s.ModName,
		Pkg:     hdlPkg,
		Handler: hdlStruct,
	}, handlerTemp)
	if err != nil {
		return
	}

	hdlFile := fmt.Sprintf("%s/%s.go", hdlDir, hdlPkg)
	err = utils.NotExistCreateFile(hdlFile, hdl)
	if err != nil {
		return
	}

	return
}
