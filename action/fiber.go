package action

import (
	"errors"
	"fmt"
	"github.com/tdeken/grass/basic"
	"github.com/tdeken/grass/createswagger"
	"github.com/tdeken/grass/utils"
	"gopkg.in/yaml.v3"
	"os"
)

type Fiber struct {
	basic.Basic
	moduleName string
}

func (s *Fiber) Error() (err error) {
	if s.Err == nil {
		return
	}

	return errors.New(fmt.Sprintf("Fiber err: %v", s.Err))
}

// NewFiber create_project
func NewFiber(rootDir string, moduleName string) *Fiber {
	fa := new(Fiber)
	{
		fa.Init(rootDir)
	}

	fa.moduleName = moduleName

	return fa
}

func (s *Fiber) Run() {
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

	err = s.meet()
	if err != nil {
		return
	}

	err = s.service()
	if err != nil {
		return
	}

	err = s.handler()
	if err != nil {
		return
	}

	err = s.route()
	if err != nil {
		return
	}

	err = s.swagger()
	if err != nil {
		return
	}

	return
}

func (s *Fiber) route() (err error) {
	path := s.PrefixDir(s.Conf.Analyze.Handler)

	entrys, err := os.ReadDir(path)
	if err != nil {
		return
	}

	var er = RouteTemp{
		Pkgs:    nil,
		Modules: nil,
	}
	for _, v := range entrys {
		if !v.IsDir() {
			continue
		}

		var filePath = s.PrefixDir(fmt.Sprintf("%s/%s", s.Conf.Analyze.Handler, v.Name()))
		var exit bool
		exit, err = utils.IsFileExist(filePath + "/controller.gen.go")
		if err != nil {
			return
		}

		if exit {
			pkg := fmt.Sprintf("%s/%s/%s", s.Conf.ModName, s.Conf.Analyze.Handler, v.Name())

			er.Modules = append(er.Modules, v.Name())
			er.Pkgs = append(er.Pkgs, pkg)
		}
	}

	text, err := utils.CreateTmp(er, routeTemp)
	if err != nil {
		return
	}

	err = utils.CreateFile(s.PrefixDir("internal/fiber/route/route.gen.go"), text)
	if err != nil {
		return
	}

	return
}

func (s *Fiber) loadConf() (err error) {
	cfd := s.PrefixDir(fmt.Sprintf("etc/grass.yaml"))

	in, err := os.ReadFile(cfd)
	if err != nil {
		err = errors.New(fmt.Sprintf("读取grass.yaml失败: %v", err))
		return
	}

	err = yaml.Unmarshal(in, &s.Conf)
	if err != nil {
		err = errors.New(fmt.Sprintf("读取grass.yaml参数异常: %v", err))
		return
	}

	protoDir := s.PrefixDir(fmt.Sprintf("%s/%s/%s.yaml", s.Conf.Proto.Path, s.moduleName, s.moduleName))
	exist, _ := utils.IsFileExist(protoDir)
	if !exist {
		err = errors.New(fmt.Sprintf("找不到当前模块的%s文件，请通过 -bp 创建", protoDir))
		return
	}

	return
}

func (s *Fiber) dir() (err error) {
	serDir, hdlDir, srsDir := s.PrefixDir(s.Conf.Analyze.Service), s.PrefixDir(s.Conf.Analyze.Handler), s.PrefixDir(s.Conf.Analyze.Sources)

	err = utils.MkDirAll(serDir, hdlDir, srsDir)
	if err != nil {
		return
	}

	serPkg, serStruct := utils.PkgAndStruct(s.Conf.Analyze.Service)
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

	hdlPkg, hdlStruct := utils.PkgAndStruct(s.Conf.Analyze.Handler)
	hdl, err := utils.CreateTmp(HandlerTemp{
		ModName: s.Conf.ModName,
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

func (s *Fiber) meet() (err error) {
	return newCreateMeet(s.Basic, s.moduleName).run()
}

func (s *Fiber) service() (err error) {
	return newCreateService(s.Basic, s.moduleName).run()
}

func (s *Fiber) handler() (err error) {
	return newCreateHandler(s.Basic, s.moduleName, basic.Fiber).run()
}

func (s *Fiber) swagger() (err error) {
	swagger := createswagger.New(s.Dir, s.moduleName)
	swagger.Run()

	return swagger.Error()
}
