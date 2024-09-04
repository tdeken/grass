package action

import (
	"bytes"
	"fmt"
	"github.com/tdeken/grass/basic"
	"github.com/tdeken/grass/utils"
	"gopkg.in/yaml.v3"
	"html/template"
	"io"
	"os"
	"reflect"
	"strings"
)

type createHandler struct {
	basic.Basic
	protoModuleName string
	routePrefix     string
	webMark         string
}

func newCreateHandler(basic basic.Basic, protoModuleName, webMark string) *createHandler {
	return &createHandler{
		Basic:           basic,
		protoModuleName: protoModuleName,
		webMark:         webMark,
	}
}

func (s *createHandler) run() (err error) {
	err = s.setRoutePrefix()
	if err != nil {
		return
	}

	parses, err := s.LoadParses(s.protoModuleName)
	if err != nil {
		return
	}

	err = s.controllers(parses)
	if err != nil {
		return
	}

	err = s.file(parses)
	if err != nil {
		return
	}

	return
}

func (s *createHandler) setRoutePrefix() (err error) {
	path := s.PrefixDir(fmt.Sprintf("%s/%s/%s.yaml", s.Conf.Proto.Path, s.protoModuleName, s.protoModuleName))

	file, err := os.OpenFile(path, os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		return
	}

	in, err := io.ReadAll(file)
	if err != nil {
		return
	}

	var doc basic.Doc
	err = yaml.Unmarshal(in, &doc)
	if err != nil {
		return
	}

	s.routePrefix = doc.Route
	return
}

func (s *createHandler) controllers(parses []basic.Parse) (err error) {
	path := s.PrefixDir(fmt.Sprintf("%s/%s", s.Conf.Analyze.Handler, s.protoModuleName))
	err = utils.NotExistCreateDir(path)
	if err != nil {
		return
	}

	var controlGen = HandlerControllerGenTemp{
		ModuleName:    s.protoModuleName,
		ModuleRoute:   s.routePrefix,
		HasController: false,
		Controllers:   nil,
		Messages:      nil,
		ModName:       s.Conf.ModName,
		ServicePath:   s.Conf.Analyze.Service,
	}

	for _, v := range parses {
		var msg string
		msg, err = s._group(v)
		if err != nil {
			return
		}

		controlGen.Messages = append(controlGen.Messages, template.HTML(msg))

		if v.Group.NotRegister {
			continue
		}
		controlGen.Controllers = append(controlGen.Controllers, v.Group.Name)
	}

	controlGen.HasController = len(controlGen.Controllers) > 0

	text, err := utils.CreateTmp(controlGen, s.getHandlerControllerGenTemp())
	if err != nil {
		return
	}

	filename := s.PrefixDir(fmt.Sprintf("%s/%s/controller.gen.go", s.Conf.Analyze.Handler, s.protoModuleName))

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
	defer file.Close()
	if err != nil {
		return
	}

	_, err = file.WriteString(text)
	if err != nil {
		return
	}

	err = s.Gofmt(filename)
	if err != nil {
		return
	}

	hg := s.Conf.Analyze.Handler[strings.LastIndex(s.Conf.Analyze.Handler, "/")+1:]
	var control = HandlerControllerTemp{
		ModName:     s.Conf.ModName,
		ModuleName:  s.protoModuleName,
		HandlerPath: s.Conf.Analyze.Handler,
		HandlerPkg:  hg,
		HandlerName: strings.ToUpper(hg[:1]) + hg[1:],
	}

	text, err = utils.CreateTmp(control, s.getHandlerControllerTemp())
	if err != nil {
		return
	}

	controlFileName := s.PrefixDir(fmt.Sprintf("%s/%s/controller.go", s.Conf.Analyze.Handler, s.protoModuleName))
	err = utils.NotExistCreateFile(controlFileName, text)
	if err != nil {
		return
	}

	err = s.Gofmt(controlFileName)
	if err != nil {
		return
	}

	return
}

func (s *createHandler) _group(parse basic.Parse) (text string, err error) {
	var group = HandlerGroupTemp{
		Name:       parse.Group.Name,
		Desc:       parse.Group.Desc,
		Route:      utils.MidString(parse.Group.Name, '-'),
		Actions:    nil,
		ModuleName: s.protoModuleName,
	}

	if parse.Group.As != "" {
		group.Route = utils.MidString(parse.Group.As, '-')
	}

	if parse.Group.NotRoute {
		group.Route = ""
	}

	for _, v := range parse.Interfaces {
		var method = "GET"
		var useLastPath, midType string
		if v.LastPath != "" {
			useLastPath = fmt.Sprintf(", action.UseLastPath(\"%s\")", v.LastPath)
		}

		if v.MidType != nil {
			var midTypeVal string
			switch reflect.TypeOf(v.MidType).Kind() {
			case reflect.String:
				midTypeVal = fmt.Sprintf("\"%v\"", v.MidType)
			default:
				midTypeVal = fmt.Sprintf("%v", v.MidType)
			}

			midType = fmt.Sprintf(", action.UseMidType(%s)", midTypeVal)
		}

		if v.Method != "" {
			method = strings.ToUpper(v.Method)
		}

		group.Actions = append(group.Actions, template.HTML(fmt.Sprintf("\t\taction.NewAction(\"%s\", s.%s%s%s),", method, v.Name, useLastPath, midType)))
	}

	return utils.CreateTmp(group, handlerGroupTemp)
}

func (s *createHandler) file(parses []basic.Parse) (err error) {
	var fileTmp = HandlerFileTemp{
		ModuleName: s.protoModuleName,
		ModName:    s.Conf.ModName,
		ParamsPath: s.Conf.Analyze.Sources,
	}

	text, err := utils.CreateTmp(fileTmp, s.getHandlerFileTemp())
	if err != nil {
		return
	}

	for _, v := range parses {
		var filename = s.PrefixDir(fmt.Sprintf("%s/%s/%s.go", s.Conf.Analyze.Handler, s.protoModuleName, utils.MidString(v.Group.Name, '_')))
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
			if strings.Contains(str.String(), fmt.Sprintf("%s(ctx *fiber.Ctx)", v1.Name)) {
				continue
			}

			var lastPath = v1.LastPath
			if lastPath == "" {
				lastPath = utils.MidString(v1.Name, '-')
			}

			var uri string
			//if len(v1.Uri) > 0 {
			//	for _, v2 := range v1.Uri {
			//		name := utils.MidString(v2.Name, '_')
			//		uri += fmt.Sprintf("/:%s", name)
			//	}
			//}

			var route string
			if s.routePrefix != "" {
				route += "/" + s.routePrefix
			}
			if !v.Group.NotRoute {
				if v.Group.As != "" {
					route += "/" + utils.MidString(v.Group.As, '-')
				} else {
					route += "/" + utils.MidString(v.Group.Name, '-')
				}
			}

			route += fmt.Sprintf("/%s%s", lastPath, uri)

			var apdTmp = HandlerFuncTemp{
				Name:   v1.Name,
				Desc:   strings.TrimPrefix(v1.Desc, "-"),
				Group:  v.Group.Name,
				Req:    fmt.Sprintf("%s%sReq", v.Group.Name, v1.Name),
				Route:  route,
				Method: strings.ToUpper(v1.Method),
			}

			var apd string
			apd, err = utils.CreateTmp(apdTmp, s.getHandlerFuncTemp())
			if err != nil {
				return
			}

			_, err = file.WriteString(apd)
			if err != nil {
				return
			}

			err = s.Gofmt(filename)
			if err != nil {
				return
			}
		}
	}

	return
}

func (s *createHandler) getHandlerControllerGenTemp() string {
	switch s.webMark {
	case basic.Gin:
		return ginHandlerControllerGenTemp
	case basic.Fiber:
		return handlerControllerGenTemp
	}
	return ""
}

func (s *createHandler) getHandlerControllerTemp() string {
	switch s.webMark {
	case basic.Gin:
		return ginHandlerControllerTemp
	case basic.Fiber:
		return handlerControllerTemp
	}
	return ""
}

func (s *createHandler) getHandlerFuncTemp() string {
	switch s.webMark {
	case basic.Gin:
		return ginHandlerFuncTemp
	case basic.Fiber:
		return handlerFuncTemp
	}
	return ""
}

func (s *createHandler) getHandlerFileTemp() string {
	switch s.webMark {
	case basic.Gin:
		return ginHandlerFileTemp
	case basic.Fiber:
		return handlerFileTemp
	}
	return ""
}
