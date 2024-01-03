package createproject

import (
	"errors"
	"fmt"
	"github.com/tdeken/grass/basic"
	"github.com/tdeken/grass/utils"
	"os/exec"
)

type CreateProject struct {
	basic.Basic
}

// NewCreateProject create_project
func NewCreateProject(modName string) *CreateProject {
	cp := new(CreateProject)
	{
		cp.Init(modName)
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

	cmd := exec.Command("go", "mod", "init", s.ModName)
	cmd.Dir = s.Dir

	_ = cmd.Run()

	return
}

func (s *CreateProject) Error() error {
	if s.Err == nil {
		return nil
	}
	return errors.New(fmt.Sprintf("CreateProject err: %v", s.Err))
}
