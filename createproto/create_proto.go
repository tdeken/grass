package createproto

import (
	"errors"
	"fmt"
	"github.com/tdeken/grass/basic"
	"github.com/tdeken/grass/utils"
	"gopkg.in/yaml.v3"
	"os"
)

type CreateProto struct {
	basic.Basic
	conf       basic.GrassConf
	moduleName string
}

func (s *CreateProto) Error() error {
	if s.Err == nil {
		return nil
	}

	return errors.New(fmt.Sprintf("CreateProto err: %v", s.Err))
}

// NewCreateProto create proto
func NewCreateProto(dir, proto string) *CreateProto {
	cp := &CreateProto{moduleName: proto}
	cp.Init(dir)

	return cp
}

// Run do logic
func (s *CreateProto) Run() {
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

func (s *CreateProto) loadConf() (err error) {
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

func (s *CreateProto) dir() (err error) {
	protoDir := s.PrefixDir(s.conf.Proto.Path)
	mDir := protoDir + "/" + s.moduleName

	err = utils.MkDirAll(mDir)
	if err != nil {
		return
	}

	return
}