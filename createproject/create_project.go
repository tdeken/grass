package createproject

import (
	"errors"
	"fmt"
	"github.com/tdeken/grass/utils"
	"os/exec"
	"strings"
)

type CreateProject struct {
	dir     string
	modName string
	err     error
}

// NewCreateProject create_project
func NewCreateProject(name string) *CreateProject {
	return &CreateProject{
		dir:     name[strings.LastIndex(name, "/")+1:],
		modName: name,
		err:     nil,
	}
}

// Run do logic
func (s *CreateProject) Run() {
	var err error
	defer func() {
		s.err = err
	}()
	err = utils.NotExistCreateDir(s.dir)
	if err != nil {
		return
	}

	cmd := exec.Command("go", "mod", "init", s.modName)
	cmd.Dir = s.dir

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
