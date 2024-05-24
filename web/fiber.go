package web

import (
	"errors"
	"fmt"
	"github.com/tdeken/grass/basic"
	"github.com/tdeken/grass/createproject"
	"github.com/tdeken/grass/utils"
	"strings"
)

type Fiber struct {
	basic.Basic
	dir     string
	modName string
}

// NewFiber create_project
func NewFiber(modName, dir string) *Fiber {

	if dir == "" {
		dir = modName
	}

	fw := new(Fiber)
	{
		fw.Init(dir)
	}

	fw.dir = fmt.Sprintf("%s/internal/fiber", fw.Dir)
	fw.modName = modName
	if fw.modName == "" {
		fw.modName = fw.dir[strings.LastIndex(fw.dir, "/")+1:]
	}

	return fw
}

func (s *Fiber) Error() error {
	if s.Err == nil {
		return nil
	}

	return errors.New(fmt.Sprintf("Fiber err: %v", s.Err))
}

// Run do logic
func (s *Fiber) Run() {
	var err error
	defer func() {
		s.Err = err
	}()
	err = s.createProject()
	if err != nil {
		return
	}

	err = s.fiber()
	if err != nil {
		return
	}

	err = s.code()
	if err != nil {
		return
	}

	err = s.config()
	if err != nil {
		return
	}

	err = s.boot()
	if err != nil {
		return
	}

	err = s.main()
	if err != nil {
		return
	}

	err = utils.RunCommand(s.Dir, "go", "mod", "tidy")
	err = utils.RunCommand(s.Dir, "go", "get", "github.com/tdeken/fiberaction@v0.1.0")
	if err != nil {
		return
	}

	return
}

// createProject init project env
func (s *Fiber) createProject() error {
	cp := createproject.NewCreateProject(s.modName, s.Dir)
	cp.Run()
	return cp.Error()
}

// fiber build fiber dir
func (s *Fiber) fiber() (err error) {
	err = utils.MkDirAll(
		s.dir,
		s.dir+"/server",
		s.dir+"/route",
		s.dir+"/result",
		s.dir+"/validate",
	)

	if err != nil {
		return
	}

	err = s.server()
	if err != nil {
		return
	}

	err = s.validate()
	if err != nil {
		return
	}

	err = s.route()
	if err != nil {
		return
	}

	err = s.result()
	if err != nil {
		return
	}

	err = s.fiberMain()
	if err != nil {
		return
	}

	return
}

// server fiber run
func (s *Fiber) server() (err error) {
	content, err := utils.CreateTmp(nil, fiberServerTemp)
	if err != nil {
		return
	}

	filename := s.dir + "/server/server.go"
	err = utils.NotExistCreateFile(filename, content)
	if err != nil {
		return
	}

	return
}

// validate form validate
func (s *Fiber) validate() (err error) {
	content, err := utils.CreateTmp(nil, fiberValidateTemp)
	if err != nil {
		return
	}

	filename := s.dir + "/validate/validate.go"
	err = utils.NotExistCreateFile(filename, content)
	if err != nil {
		return
	}

	return
}

// route load route
func (s *Fiber) route() (err error) {
	content, err := utils.CreateTmp(nil, routeTemp)
	if err != nil {
		return
	}

	filename := s.dir + "/route/route.gen.go"
	err = utils.NotExistCreateFile(filename, content)
	if err != nil {
		return
	}

	return
}

// code error code
func (s *Fiber) code() (err error) {
	dir := fmt.Sprintf("%s/internal/code", s.Dir)

	err = utils.NotExistCreateDir(dir)
	if err != nil {
		return
	}

	content, err := utils.CreateTmp(nil, codeTemp)
	if err != nil {
		return
	}

	filename := dir + "/code.go"
	err = utils.NotExistCreateFile(filename, content)
	if err != nil {
		return
	}

	content, err = utils.CreateTmp(nil, codeErrorTemp)
	if err != nil {
		return
	}

	filename = dir + "/error.go"
	err = utils.NotExistCreateFile(filename, content)
	if err != nil {
		return
	}

	return
}

// result response data struct
func (s *Fiber) result() (err error) {
	content, err := utils.CreateTmp(FiberResultTemp{ModName: s.modName}, fiberResultTemp)
	if err != nil {
		return
	}

	filename := s.dir + "/result/result.go"
	err = utils.NotExistCreateFile(filename, content)
	if err != nil {
		return
	}

	return
}

// fiber response data struct
func (s *Fiber) fiberMain() (err error) {
	content, err := utils.CreateTmp(FiberTemp{ModName: s.modName}, fiberTemp)
	if err != nil {
		return
	}

	filename := s.dir + "/fiber.go"
	err = utils.NotExistCreateFile(filename, content)
	if err != nil {
		return
	}

	return
}

// config load config
func (s *Fiber) config() (err error) {
	dir := s.Dir + "/internal/config"

	err = utils.NotExistCreateDir(dir)
	if err != nil {
		return
	}

	content, err := utils.CreateTmp(nil, configTemp)
	if err != nil {
		return
	}

	filename := dir + "/config.go"
	err = utils.NotExistCreateFile(filename, content)
	if err != nil {
		return
	}

	var env = []string{"local", "dev", "test", "prod"}
	for _, v := range env {
		content, err = utils.CreateTmp(nil, configFileTemp)
		if err != nil {
			return
		}

		filename = fmt.Sprintf("%s/etc/config-%s.yaml", s.Dir, v)
		err = utils.NotExistCreateFile(filename, content)
		if err != nil {
			return
		}
	}

	return
}

// boot load boot
func (s *Fiber) boot() (err error) {
	dir := s.Dir + "/internal/boot"

	err = utils.NotExistCreateDir(dir)
	if err != nil {
		return
	}

	content, err := utils.CreateTmp(FiberBootTemp{ModName: s.modName}, fiberBootTemp)
	if err != nil {
		return
	}

	filename := dir + "/boot.go"
	err = utils.NotExistCreateFile(filename, content)
	if err != nil {
		return
	}

	return
}

// main
func (s *Fiber) main() (err error) {

	content, err := utils.CreateTmp(MainTemp{
		ModName: s.modName,
		Spc:     "<-",
	}, mainTemp)
	if err != nil {
		return
	}

	filename := s.Dir + "/main.go"
	err = utils.NotExistCreateFile(filename, content)
	if err != nil {
		return
	}

	content, err = utils.CreateTmp(DockerfileTemp{ModName: s.modName}, dockerfileTemp)
	if err != nil {
		return
	}

	filename = s.Dir + "/Dockerfile"
	err = utils.NotExistCreateFile(filename, content)
	if err != nil {
		return
	}

	return
}
