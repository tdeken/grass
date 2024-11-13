package createswagger

import (
	"errors"
	"fmt"
	"github.com/tdeken/grass/basic"
	"github.com/tdeken/grass/utils"
	"gopkg.in/yaml.v3"
	"html/template"
	"io"
	"os"
	"strings"
)

type CreateSwagger struct {
	basic.Basic
	moduleName           string
	route, securityTitle string
	dictionary           map[string]string
}

func (s *CreateSwagger) Error() error {
	if s.Err == nil {
		return nil
	}

	return errors.New(fmt.Sprintf("CreateSwagger err: %v", s.Err))
}

// New create swagger
func New(dir, moduleName string) *CreateSwagger {
	cp := &CreateSwagger{moduleName: moduleName}
	cp.Init(dir)

	return cp
}

// Run do logic
func (s *CreateSwagger) Run() {
	var err error
	defer func() {
		s.Err = err
	}()

	err = s.LoadConf()
	if err != nil {
		return
	}

	if s.Conf.Swagger.Path == "" {
		return
	}

	err = utils.MkDirAll(s.PrefixDir(s.Conf.Swagger.Path + "/" + s.moduleName))
	if err != nil {
		return
	}

	err = s.doc()
	if err != nil {
		return
	}

	err = s.file()
	if err != nil {
		return
	}

	err = s.generate()
	if err != nil {
		return
	}

	return
}

func (s *CreateSwagger) doc() (err error) {
	file, err := os.OpenFile(s.PrefixDir(fmt.Sprintf("%s/%s/%s.yaml", s.Conf.Proto.Path, s.moduleName, s.moduleName)), os.O_RDWR, 0777)
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

	s.route = doc.Route
	s.securityTitle = doc.Auth.Security

	var desc string

	if doc.Title != "" {
		desc += "\n" + "// @title " + doc.Title
	}

	if doc.Host != "" {
		desc += "\n" + "// @host " + doc.Host
	}

	if doc.Schemes != nil {
		desc += "\n" + "// @schemes " + strings.Join(doc.Schemes, " ")
	}

	if doc.Ver != "" {
		desc += "\n" + "// @version " + doc.Ver
	}

	if doc.Desc != "" {
		desc += "\n" + "// @description " + doc.Desc
	}

	if doc.Auth.Security != "" {
		desc += "\n" + "// @securityDefinitions." + doc.Auth.Security + " " + doc.Auth.Title
	}

	if doc.Auth.In != "" {
		desc += "\n" + "// @in " + doc.Auth.In
	}

	if doc.Auth.Name != "" {
		desc += "\n" + "// @name " + doc.Auth.Name
	}

	if doc.Contact.Name != "" {
		desc += "\n" + "// @contact.name " + doc.Contact.Name
	}

	if doc.Contact.Url != "" {
		desc += "\n" + "// @contact.url " + doc.Contact.Url
	}

	if doc.Contact.Email != "" {
		desc += "\n" + "// @contact.email " + doc.Contact.Email
	}

	//模板内容
	var content = SwaggerDocTemp{
		ModuleName: s.moduleName,
		Content:    template.HTML(desc),
	}

	text, err := utils.CreateTmp(content, swaggerDocTemp)
	if err != nil {
		return
	}

	err = utils.CreateFile(s.PrefixDir(fmt.Sprintf("%s/%s/doc.go", s.Conf.Swagger.Path, s.moduleName)), text)
	if err != nil {
		return
	}

	return
}

func (s *CreateSwagger) file() (err error) {
	parses, err := s.LoadParses(s.moduleName)
	if err != nil {
		return
	}

	for _, v := range parses {
		var filename = s.PrefixDir(fmt.Sprintf("%s/%s/%s.go", s.Conf.Swagger.Path, s.moduleName, utils.MidString(v.Group.Name, '_')))

		var content = SwaggerFileTemp{
			ModuleName: s.moduleName,
			Group:      v.Group.Name,
		}

		var text string
		text, err = utils.CreateTmp(content, swaggerFileTemp)
		if err != nil {
			return
		}

		err = utils.CreateFile(filename, text)
		if err != nil {
			return
		}

		err = s._swag(v)
		if err != nil {
			return
		}

		err = utils.RunCommand("", "gofmt", "-w", filename)
		if err != nil {
			return
		}
	}
	return
}

func (s *CreateSwagger) _swag(parse basic.Parse) (err error) {
	var filename = s.PrefixDir(fmt.Sprintf("%s/%s/%s.go", s.Conf.Swagger.Path, s.moduleName, utils.MidString(parse.Group.Name, '_')))
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND, 0777)
	defer file.Close()
	if err != nil {
		return
	}

	for _, v := range parse.Interfaces {

		var lastPath = v.LastPath
		if lastPath == "" {
			lastPath = utils.MidString(v.Name, '-')
		}

		var method = v.Method
		if method == "" {
			method = "GET"
		}

		var reqContentType = "application/json"
		if v.ReqContentType != "" {
			reqContentType = v.ReqContentType
		}

		var resContentType = "application/json"
		if v.ResContentType != "" {
			resContentType = v.ResContentType
		}

		var body = "body"
		if strings.ToUpper(method) == "GET" {
			body = "query"
		}

		var annotation []template.HTML

		var uri string

		var groupDesc = parse.Group.Desc
		if parse.Group.Tags != "" {
			groupDesc = parse.Group.Tags
		}
		var desc = v.Desc
		if strings.Index(v.Desc, "-") == 0 {
			desc = parse.Group.Desc + v.Desc
		}

		var sw = SwaggerTemp{
			Name:           v.Name,
			GroupDesc:      groupDesc,
			Desc:           desc,
			SecurityTitle:  s.securityTitle,
			ReqContentType: reqContentType,
			ResContentType: resContentType,
			Body:           body,
			Req:            parse.Group.Name + v.Name + "Req",
			ResFormat:      parse.Group.Name + v.Name + "Format",
			Route:          "",
			Method:         strings.ToUpper(method),
			Group:          parse.Group.Name,
			Messages:       nil,
			Annotation:     annotation,
		}

		if s.route != "" {
			sw.Route += "/" + s.route
		}

		if !parse.Group.NotRoute {
			if parse.Group.As != "" {
				sw.Route += "/" + utils.MidString(parse.Group.As, '-')
			} else {
				sw.Route += "/" + utils.MidString(parse.Group.Name, '-')
			}
		}

		sw.Route += fmt.Sprintf("/%s%s", lastPath, uri)

		sw.Messages = append(sw.Messages,
			template.HTML(fmt.Sprintf(
				"type %s struct {\n\t\t%s int32 `json:\"%s\"`\n\t\t%s string `json:\"%s\"`\n\t\t%s %s `json:\"%s\"`\n}",
				parse.Group.Name+v.Name+"Format",
				utils.CamelString(s.Conf.Swagger.Code),
				s.Conf.Swagger.Code,
				utils.CamelString(s.Conf.Swagger.Msg),
				s.Conf.Swagger.Msg,
				utils.CamelString(s.Conf.Swagger.Data),
				parse.Group.Name+v.Name+"Res",
				s.Conf.Swagger.Data,
			)),
		)

		prefix := parse.Group.Name + v.Name

		s.parseMsg(&sw, v.Msgs, prefix)

		v.Req.Name = prefix
		//if len(v.Uri) > 0 {
		//	s.parseUri(&sw, v.Req.Name, v.Uri, "Uri")
		//}
		s.parseReq(&sw, v.Req, "Req")
		v.Res.Name = prefix
		s.parseRes(&sw, v.Res, "Res")

		var text string
		text, err = utils.CreateTmp(sw, swaggerTemp)
		if err != nil {
			return
		}

		_, err = file.WriteString(text)
		if err != nil {
			return
		}
	}

	return
}

// 解析通用结构体
func (s *CreateSwagger) parseMsg(sw *SwaggerTemp, msgs []basic.Msg, prefix string) {
	s.dictionary = map[string]string{}

	for _, msg := range msgs {
		s.dictionary[msg.Name] = prefix + msg.Name
	}

	for _, res := range msgs {
		name := prefix + res.Name
		sn := "type " + name + " struct { \n"
		for _, field := range res.Fields {
			var jsonTag string
			if field.Omit {
				jsonTag = ",omitempty"
			}

			sn += fmt.Sprintf("\t%s %s `json:\"%s%s\"` //%s\n", utils.CamelString(field.Name), utils.GetClass(name, field.Class, s.dictionary), field.Name, jsonTag, field.Desc)
		}
		sn += "}"

		sw.Messages = append(sw.Messages, template.HTML(sn))
	}
}

// 解析请求结构体
func (s *CreateSwagger) parseReq(sw *SwaggerTemp, req basic.Message, suffix string) {
	sn := "type " + req.Name + suffix + " struct { \n"
	for _, field := range req.Fields {
		var bind string
		if field.Validate != "" {
			bind = " validate:\"" + field.Validate + "\""
		}

		sn += fmt.Sprintf("\t%s %s `json:\"%s\" %s` //%s\n", utils.CamelString(field.Name), utils.GetClass(req.Name, field.Class, s.dictionary), field.Name, bind, field.Desc)
	}
	sn += "}"

	sw.Messages = append(sw.Messages, template.HTML(sn))

	for _, v := range req.Messages {
		v.Name = req.Name + v.Name
		s.parseReq(sw, v, "")
	}
}

// 解析返回结构体
func (s *CreateSwagger) parseRes(sw *SwaggerTemp, res basic.Message, suffix string) {
	sn := "type " + res.Name + suffix + " struct { \n"
	for _, field := range res.Fields {
		var jsonTag string
		if field.Omit {
			jsonTag = ",omitempty"
		}
		sn += fmt.Sprintf("\t%s %s `json:\"%s%s\"` //%s\n", utils.CamelString(field.Name), utils.GetClass(res.Name, field.Class, s.dictionary), field.Name, jsonTag, field.Desc)
	}
	sn += "}"

	sw.Messages = append(sw.Messages, template.HTML(sn))

	for _, v := range res.Messages {
		v.Name = res.Name + v.Name
		s.parseRes(sw, v, "")
	}
}

func (s *CreateSwagger) generate() (err error) {
	filename := s.PrefixDir("generate.go")
	err = utils.NotExistCreateFile(filename, "package main\n\n")
	if err != nil {
		return
	}

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND, 0777)
	defer file.Close()
	if err != nil {
		return
	}

	b, err := io.ReadAll(file)
	if err != nil {
		return
	}

	var content = fmt.Sprintf("//go:generate swag init  -o docs/%s -g doc.go -d %s/%s", s.moduleName, s.Conf.Swagger.Path, s.moduleName)
	text := string(b)

	if !strings.Contains(text, "//go:generate swag init") {
		err = utils.RunCommand(s.Dir, "go", "get", "github.com/swaggo/swag")
		if err != nil {
			return
		}
	}

	if strings.Contains(text, content) {
		return
	}

	_, err = file.WriteString(content + "\r\n")
	return
}
