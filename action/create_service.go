package action

import (
	"bytes"
	"fmt"
	"github.com/tdeken/grass/basic"
	"github.com/tdeken/grass/utils"
	"io"
	"os"
	"strings"
)

type createService struct {
	basic.Basic
	protoModuleName string
}

func newCreateService(basic basic.Basic, protoModuleName string) *createService {
	return &createService{
		Basic:           basic,
		protoModuleName: protoModuleName,
	}
}

func (s *createService) run() (err error) {
	path := s.PrefixDir(fmt.Sprintf("%s/%s", s.Conf.Analyze.Service, s.protoModuleName))

	err = utils.NotExistCreateDir(path)
	if err != nil {
		return
	}
	return s.file()
}

func (s *createService) file() (err error) {
	path := s.PrefixDir(fmt.Sprintf("%s/%s", s.Conf.Analyze.Service, s.protoModuleName))

	var fileTmp = ServiceGroupTemp{
		ModName:     s.Conf.ModName,
		ModuleName:  s.protoModuleName,
		ServicePath: s.Conf.Analyze.Service,
		ServicePkg:  s.Conf.Analyze.Service[strings.LastIndex(s.Conf.Analyze.Service, "/")+1:],
		Name:        "",
		ParamsPath:  fmt.Sprintf("%s/%s", s.Conf.Analyze.Sources, s.protoModuleName),
	}

	parses, err := s.LoadParses(s.protoModuleName)
	if err != nil {
		return
	}

	for _, v := range parses {
		var filename = fmt.Sprintf("%s/%s.go", path, utils.MidString(v.Group.Name, '_'))
		fileTmp.Name = v.Group.Name
		var text string
		text, err = utils.CreateTmp(fileTmp, serviceGroupTemp)
		if err != nil {
			return
		}

		err = utils.NotExistCreateFile(filename, text)
		if err != nil {
			return
		}

		for _, v1 := range v.Interfaces {
			var file *os.File

			file, err = os.OpenFile(filename, os.O_APPEND|os.O_RDWR, 0777)
			if err != nil {
				return
			}

			var b []byte
			b, err = io.ReadAll(file)
			if err != nil {
				return
			}
			var str = bytes.NewBuffer(b)
			if strings.Contains(str.String(), fmt.Sprintf("func (s %s) %s(", v.Group.Name, v1.Name)) {
				continue
			}

			var apdTmp = ServiceFuncTemp{
				Name:      v1.Name,
				Desc:      strings.TrimPrefix(v1.Desc, "-"),
				GroupName: v.Group.Name,
			}

			var apd string
			apd, err = utils.CreateTmp(apdTmp, serviceFuncTemp)
			if err != nil {
				return
			}

			_, err = file.WriteString(apd)
			if err != nil {
				return
			}
			file.Close()

			err = s.Gofmt(filename)
			if err != nil {
				return
			}
		}

	}
	return
}
