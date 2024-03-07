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
	temp       bool
}

func (s *CreateProto) Error() error {
	if s.Err == nil {
		return nil
	}

	return errors.New(fmt.Sprintf("CreateProto err: %v", s.Err))
}

// NewCreateProto create proto
func NewCreateProto(dir, proto string, temp bool) *CreateProto {
	cp := &CreateProto{moduleName: proto, temp: temp}
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

	content, err := utils.CreateTmp(ProtoTemp{
		Route: s.moduleName,
	}, protoTemp)
	if err != nil {
		return
	}

	filename := fmt.Sprintf("%s/%s.yaml", mDir, s.moduleName)
	err = utils.NotExistCreateFile(filename, content)
	if err != nil || !s.temp {
		return
	}

	var temp, tempFile string

	switch s.conf.Proto.FileType {
	case "json":
		temp, err = utils.CreateTmp(nil, exampleTemp)
		if err != nil {
			return
		}
		tempFile = fmt.Sprintf("%s/example.json", mDir)
	case "toml":
		temp, err = utils.CreateTmp(nil, exampleTomlTemp)
		if err != nil {
			return
		}
		tempFile = fmt.Sprintf("%s/example.toml", mDir)
	}

	err = utils.NotExistCreateFile(tempFile, temp)
	if err != nil {
		return
	}

	return
}
