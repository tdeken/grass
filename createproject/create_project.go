package createproject

import (
	"errors"
	"fmt"
	"grass/utils"
	"os/exec"
)

type CreateProject struct {
	Name string
	err  error
}

func NewCreateProject(name string) *CreateProject {
	return &CreateProject{
		Name: name,
		err:  nil,
	}
}

func (s *CreateProject) Run() {
	var err error
	defer func() {
		s.err = err
	}()
	err = utils.NotExistCreateDir(s.Name)
	if err != nil {
		return
	}

	cmd := exec.Command("go", "mod", "init", s.Name)
	cmd.Dir = s.Name

	err = cmd.Run()
	if err != nil {
		return
	}

	return
}

func (s *CreateProject) Error() error {
	if s.err == nil {
		return nil
	}
	return s.errPrefix(s.err)
}

func (s *CreateProject) errPrefix(err error) error {
	return errors.New(fmt.Sprintf("CreateProject err: %v", err))
}
