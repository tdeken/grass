package fiberweb

import (
	"errors"
	"fmt"
	"github.com/tdeken/grass/basic"
	"github.com/tdeken/grass/createproject"
)

type FiberWeb struct {
	basic.Basic
}

// NewFiberWeb create_project
func NewFiberWeb(modName string) *FiberWeb {
	fw := new(FiberWeb)
	{
		fw.Init(modName)
	}
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

	return
}

func (s *FiberWeb) createProject() error {
	cp := createproject.NewCreateProject(s.ModName)
	cp.Run()
	return cp.Error()
}
