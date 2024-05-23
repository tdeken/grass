package basic

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

type Basic struct {
	Dir  string
	Err  error
	Conf *GrassConf
}

func (s *Basic) Init(dirPath string) {
	s.Dir = dirPath
	s.Err = nil
}

func (s *Basic) PrefixDir(path string) string {
	if s.Dir != "" {
		return s.Dir + "/" + path
	}
	return path
}

func (s *Basic) LoadConf() (err error) {
	cfd := s.PrefixDir(fmt.Sprintf("etc/grass.yaml"))

	in, err := os.ReadFile(cfd)
	if err != nil {
		err = errors.New(fmt.Sprintf("读取grass.yaml失败: %v", err))
		return
	}

	err = yaml.Unmarshal(in, &s.Conf)
	if err != nil {
		err = errors.New(fmt.Sprintf("读取grass.yaml参数异常: %v", err))
		return
	}

	return
}

func (s *Basic) LoadParses(protoModuleName string) (parses []Parse, err error) {
	path := s.PrefixDir(fmt.Sprintf("%s/%s", s.Conf.Proto.Path, protoModuleName))

	entrys, err := os.ReadDir(path)
	if err != nil {
		return
	}

	parses = make([]Parse, 0, len(entrys)-1)
	for _, entry := range entrys {
		if !strings.HasSuffix(entry.Name(), fmt.Sprintf(".%s", s.Conf.Proto.FileType)) {
			continue
		}

		filename := path + "/" + entry.Name()
		var parse = Parse{}
		err = parse.Parse(filename, s.Conf.Proto.FileType)
		if err != nil {
			return
		}

		parses = append(parses, parse)
	}

	return
}
