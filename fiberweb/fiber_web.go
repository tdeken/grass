package fiberweb

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/tdeken/grass/basic"
	"github.com/tdeken/grass/createproject"
	"github.com/tdeken/grass/utils"
	"os/exec"
)

type FiberWeb struct {
	basic.Basic
	fiberDir string
}

// NewFiberWeb create_project
func NewFiberWeb(modName string) *FiberWeb {
	fw := new(FiberWeb)
	{
		fw.Init(modName)
	}

	fw.fiberDir = fmt.Sprintf("%s/internal/fiber", fw.Dir)

	return fw
}

func (s *FiberWeb) Error() error {
	if s.Err == nil {
		return nil
	}

	return errors.New(fmt.Sprintf("FiberWeb err: %v", s.Err))
}

// Run do logic
func (s *FiberWeb) Run() {
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

	var out bytes.Buffer
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = s.Dir
	cmd.Stderr = &out

	err = cmd.Run()
	if err != nil {
		err = errors.New(out.String())
		return
	}

	return
}

// createProject init project env
func (s *FiberWeb) createProject() error {
	cp := createproject.NewCreateProject(s.ModName)
	cp.Run()
	return cp.Error()
}

// fiber build fiber dir
func (s *FiberWeb) fiber() (err error) {
	err = utils.MkDirAll(
		s.fiberDir,
		s.fiberDir+"/server",
		s.fiberDir+"/route",
		s.fiberDir+"/result",
		s.fiberDir+"/validate",
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
func (s *FiberWeb) server() (err error) {
	content, err := utils.CreateTmp(nil, fiberServerTemp)
	if err != nil {
		return
	}

	filename := s.fiberDir + "/server/server.go"
	err = utils.NotExistCreateFile(filename, content)
	if err != nil {
		return
	}

	return
}

// validate form validate
func (s *FiberWeb) validate() (err error) {
	content, err := utils.CreateTmp(nil, fiberValidateTemp)
	if err != nil {
		return
	}

	filename := s.fiberDir + "/validate/validate.go"
	err = utils.NotExistCreateFile(filename, content)
	if err != nil {
		return
	}

	return
}

// route load route
func (s *FiberWeb) route() (err error) {
	content, err := utils.CreateTmp(nil, fiberRouteTemp)
	if err != nil {
		return
	}

	filename := s.fiberDir + "/route/route.go"
	err = utils.NotExistCreateFile(filename, content)
	if err != nil {
		return
	}

	return
}

// code error code
func (s *FiberWeb) code() (err error) {
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
func (s *FiberWeb) result() (err error) {
	content, err := utils.CreateTmp(FiberResultTemp{ModName: s.ModName}, fiberResultTemp)
	if err != nil {
		return
	}

	filename := s.fiberDir + "/result/result.go"
	err = utils.NotExistCreateFile(filename, content)
	if err != nil {
		return
	}

	return
}

// fiber response data struct
func (s *FiberWeb) fiberMain() (err error) {
	content, err := utils.CreateTmp(FiberTemp{ModName: s.ModName}, fiberTemp)
	if err != nil {
		return
	}

	filename := s.fiberDir + "/fiber.go"
	err = utils.NotExistCreateFile(filename, content)
	if err != nil {
		return
	}

	return
}

// config load config
func (s *FiberWeb) config() (err error) {
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
func (s *FiberWeb) boot() (err error) {
	dir := s.Dir + "/internal/boot"

	err = utils.NotExistCreateDir(dir)
	if err != nil {
		return
	}

	content, err := utils.CreateTmp(BootTemp{ModName: s.ModName}, bootTemp)
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
func (s *FiberWeb) main() (err error) {

	content, err := utils.CreateTmp(MainTemp{
		ModName: s.ModName,
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

	content, err = utils.CreateTmp(DockerfileTemp{ModName: s.ModName}, dockerfileTemp)
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
