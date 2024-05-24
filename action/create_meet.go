package action

import (
	"bytes"
	"fmt"
	"github.com/tdeken/grass/basic"
	"github.com/tdeken/grass/utils"
	"html/template"
	"os"
	"strings"
)

type createMeet struct {
	basic.Basic
	protoModuleName string
}

func newCreateMeet(basic basic.Basic, protoModuleName string) *createMeet {
	return &createMeet{
		Basic:           basic,
		protoModuleName: protoModuleName,
	}
}

func (s *createMeet) run() (err error) {
	parses, err := s.LoadParses(s.protoModuleName)
	if err != nil {
		return
	}

	moduleDir := s.PrefixDir(s.Conf.Analyze.Sources) + "/" + s.protoModuleName
	err = utils.NotExistCreateDir(moduleDir)
	if err != nil {
		return
	}

	for _, parse := range parses {
		var filename string
		var meet = meetOne{MeetTemp: MeetTemp{
			ModuleName: s.protoModuleName,
			Messages:   nil,
		}}
		filename, err = meet.one(parse, s.PrefixDir(s.Conf.Analyze.Sources))
		if err != nil {
			return
		}

		err = s.Gofmt(filename)
		if err != nil {
			return
		}
	}

	return
}

type meetOne struct {
	dictionary map[string]string
	MeetTemp
}

func (b *meetOne) one(parse basic.Parse, savePath string) (filename string, err error) {
	for _, v := range parse.Interfaces {
		b.parseMsg(v.Msgs, parse.Group.Name+v.Name)

		v.Req.Name = parse.Group.Name + v.Name
		b.parseReq(v.Req, "Req", "", "", v.Method)

		v.Res.Name = parse.Group.Name + v.Name
		b.parseRes(v.Res, "Res")
	}

	t, err := template.New("pb.tpl").Parse(meetTemp)
	if err != nil {
		return
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, b)
	if err != nil {
		return
	}

	filename = fmt.Sprintf("%s/%s/%s.gen.go", savePath, b.ModuleName, utils.MidString(parse.Group.Name, '_'))

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
	defer file.Close()
	if err != nil {
		return
	}

	_, err = file.WriteString(buf.String())
	if err != nil {
		return
	}

	return
}

// 解析通用结构体
func (b *meetOne) parseMsg(msgs []basic.Msg, prefix string) {
	b.dictionary = map[string]string{}

	for _, msg := range msgs {
		b.dictionary[msg.Name] = prefix + msg.Name
	}

	for _, res := range msgs {
		name := prefix + res.Name
		sn := "type " + name + " struct { \n"
		for _, field := range res.Fields {
			var validate string
			if field.Validate != "" {
				validate = " validate:\"" + field.Validate + "\""
			}

			var jsonTag string
			if field.Omit {
				jsonTag = ",omitempty"
			}

			var tagStr string
			if field.Tags != "" {
				tags := strings.Split(field.Tags, ";")

				for _, v := range tags {
					f, s, _ := strings.Cut(v, ":")
					tagStr += fmt.Sprintf(" %s:\"%s\"", f, s)
				}
			}

			sn += fmt.Sprintf("\t%s %s `json:\"%s%s\"%s%s` //%s\n", utils.CamelString(field.Name), utils.GetClass(name, field.Class, b.dictionary), field.Name, jsonTag, validate, tagStr, field.Desc)
		}
		sn += "}"

		b.Messages = append(b.Messages, template.HTML(sn))
	}
}

// 解析请求结构体
func (b *meetOne) parseReq(req basic.Message, suffix, uriName, headerName, method string) {
	sn := "type " + req.Name + suffix + " struct { \n"
	if headerName != "" {
		sn += headerName + "\n"
	}
	if uriName != "" {
		sn += uriName + "\n"
	}

	for _, field := range req.Fields {
		var validate string
		if field.Validate != "" {
			validate = " validate:\"" + field.Validate + "\""
		}

		var jsonTag string
		if field.Omit {
			jsonTag = ",omitempty"
		}

		var from = field.From
		if strings.ToUpper(method) == "GET" && from == "" {
			from = "query"
		}

		var tagStr string
		if field.Tags != "" {
			tags := strings.Split(field.Tags, ";")

			for _, v := range tags {
				f, s, _ := strings.Cut(v, ":")
				tagStr += fmt.Sprintf(" %s:\"%s\"", f, s)
			}
		}

		sn += fmt.Sprintf("\t%s %s `json:\"%s%s\" %s%s%s` //%s\n",
			utils.CamelString(field.Name),
			utils.GetClass(req.Name, field.Class, b.dictionary),
			field.Name,
			jsonTag,
			b.parseFromTag(from, field.Name),
			validate,
			tagStr,
			field.Desc,
		)
	}
	sn += "}"

	b.Messages = append(b.Messages, template.HTML(sn))

	for _, v := range req.Messages {
		v.Name = req.Name + v.Name
		b.parseReq(v, "", "", "", method)
	}
}

func (b *meetOne) parseFromTag(from string, fieldName string) string {
	switch from {
	case "query":
		return fmt.Sprintf("query:\"%s\"", fieldName)
	default:
		return fmt.Sprintf("form:\"%s\"", fieldName)
	}
}

// 解析返回结构体
func (b *meetOne) parseRes(res basic.Message, suffix string) {
	sn := "type " + res.Name + suffix + " struct { \n"
	for _, field := range res.Fields {
		var jsonTag string
		if field.Omit {
			jsonTag = ",omitempty"
		}

		var tagStr string
		if field.Tags != "" {
			tags := strings.Split(field.Tags, ";")

			for _, v := range tags {
				f, s, _ := strings.Cut(v, ":")
				tagStr += fmt.Sprintf(" %s:\"%s\"", f, s)
			}
		}

		sn += fmt.Sprintf("\t%s %s `json:\"%s%s\"%s` //%s\n", utils.CamelString(field.Name), utils.GetClass(res.Name, field.Class, b.dictionary), field.Name, jsonTag, tagStr, field.Desc)
	}
	sn += "}"

	b.Messages = append(b.Messages, template.HTML(sn))

	for _, v := range res.Messages {
		v.Name = res.Name + v.Name
		b.parseRes(v, "")
	}
}

// 解析请求结构体
func (b *meetOne) parseUri(Name string, fields []basic.Field) {
	sn := "type " + Name + " struct { \n"
	for _, field := range fields {
		var validate string
		if field.Validate != "" {
			validate = " validate:\"" + field.Validate + "\""
		}

		var tagStr string
		if field.Tags != "" {
			tags := strings.Split(field.Tags, ";")

			for _, v := range tags {
				f, s, _ := strings.Cut(v, ":")
				tagStr += fmt.Sprintf(" %s:\"%s\"", f, s)
			}
		}

		sn += fmt.Sprintf("\t%s %s `%s%s%s` //%s\n",
			utils.CamelString(field.Name),
			field.Class,
			b.parseFromTag("uri", field.Name),
			validate,
			tagStr,
			field.Desc,
		)
	}
	sn += "}"

	b.Messages = append(b.Messages, template.HTML(sn))

}

// 解析请求结构体
func (b *meetOne) parseHeader(Name string, fields []basic.Field) {
	sn := "type " + Name + " struct { \n"
	for _, field := range fields {
		var validate string
		if field.Validate != "" {
			validate = " validate:\"" + field.Validate + "\""
		}

		var tagStr string
		if field.Tags != "" {
			tags := strings.Split(field.Tags, ";")

			for _, v := range tags {
				f, s, _ := strings.Cut(v, ":")
				tagStr += fmt.Sprintf(" %s:\"%s\"", f, s)
			}
		}

		sn += fmt.Sprintf("\t%s %s `%s%s%s` //%s\n",
			utils.CamelString(field.Name),
			field.Class,
			b.parseFromTag("header", field.Name),
			validate,
			tagStr,
			field.Desc,
		)
	}
	sn += "}"

	b.Messages = append(b.Messages, template.HTML(sn))
}
