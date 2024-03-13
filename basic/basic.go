package basic

import (
	"strings"
)

type Basic struct {
	Dir string
	Err error
}

func (s *Basic) Init(dirPath string) {
	s.Dir = dirPath[strings.LastIndex(dirPath, "/")+1:]
	s.Err = nil
}

func (s *Basic) PrefixDir(path string) string {
	if s.Dir != "" {
		return s.Dir + "/" + path
	}
	return path
}
