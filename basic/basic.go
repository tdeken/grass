package basic

import "strings"

type Basic struct {
	Dir     string
	ModName string
	Err     error
}

func (s *Basic) Init(modName string) {
	s.Dir = modName[strings.LastIndex(modName, "/")+1:]
	s.ModName = modName
	s.Err = nil
}
