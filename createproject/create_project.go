package createproject

import (
	"errors"
	"fmt"
	"github.com/tdeken/grass/basic"
	"github.com/tdeken/grass/utils"
	"strings"
)

type CreateProject struct {
	basic.Basic
}

// NewCreateProject create_project
func NewCreateProject(dirPath string) *CreateProject {
	cp := new(CreateProject)
	{
		cp.Init(dirPath)
	}

	return cp
}

// Run do logic
func (s *CreateProject) Run() {
	var err error
	defer func() {
		s.Err = err
	}()
	err = utils.NotExistCreateDir(s.Dir)
	if err != nil {
		return
	}

	err = utils.RunCommand(s.Dir, "go", "mod", "init", s.Dir[strings.LastIndex(s.Dir, "/")+1:])
	if err != nil {
		return
	}

	err = s.etc()
	if err != nil {
		return
	}

	return
}

// etc dir init
// stores configuration files
func (s *CreateProject) etc() (err error) {
	dir := s.Dir + "/etc"

	err = utils.NotExistCreateDir(dir)
	if err != nil {
		return
	}

	content, err := utils.CreateTmp(ProtoYamlFile{
		ModName: s.Dir[strings.LastIndex(s.Dir, "/")+1:],
	}, protoYamlFile)
	if err != nil {
		return
	}

	filename := dir + "/grass.yaml"
	err = utils.NotExistCreateFile(filename, content)

	return
}

func (s *CreateProject) Error() error {
	if s.Err == nil {
		return nil
	}

	return errors.New(fmt.Sprintf("CreateProject err: %v", s.Err))
}
