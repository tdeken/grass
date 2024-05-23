package createproto

import (
	"errors"
	"fmt"
	"github.com/tdeken/grass/basic"
	"github.com/tdeken/grass/utils"
)

type CreateProto struct {
	basic.Basic
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

	err = s.LoadConf()
	if err != nil {
		return
	}

	err = s.dir()
	if err != nil {
		return
	}

	return
}

func (s *CreateProto) dir() (err error) {
	protoDir := s.PrefixDir(s.Conf.Proto.Path)
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

	switch s.Conf.Proto.FileType {
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
