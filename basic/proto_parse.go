package basic

import (
	"encoding/json"
	"errors"
	"github.com/pelletier/go-toml/v2"
	"io"
	"os"
	"strings"
)

type Parse struct {
	Group      Group        `json:"group"`      //接口组名称
	Interfaces []*Interface `json:"interfaces"` //接口
}

type Group struct {
	Name        string `json:"name"`         //组名称
	As          string `json:"as"`           //路由名称
	Desc        string `json:"desc"`         //组描述
	Tags        string `json:"tags"`         //swagger分组，空的时候，默认是组描述
	NotRegister bool   `json:"not_register"` //不注册
}

type Interface struct {
	Name           string      `json:"name"`             //接口名称
	Desc           string      `json:"desc"`             //接口备注
	Method         string      `json:"method"`           //请求方法
	LastPath       string      `json:"last_path"`        //最后一节路由
	MidType        interface{} `json:"mid_type"`         //中间件类型
	ReqContentType string      `json:"req_content_type"` //接参数形式
	ResContentType string      `json:"res_content_type"` //返回参数形式
	//Header         []Field     `json:"header"`           //请求头参数
	//Uri            []Field     `json:"uri"`              //uri参数
	Req  Message `json:"req"`  //请求参数
	Res  Message `json:"res"`  //返回数据
	Msgs []Msg   `json:"msgs"` //通用结构体
}

type Field struct {
	Name     string `json:"name"`     //字段名
	Class    string `json:"class"`    //字段类型
	Desc     string `json:"desc"`     //字段备注
	Validate string `json:"validate"` //字段校验 https://github.com/go-playground/validator 语法
	From     string `json:"from"`     //字段来源
	Omit     bool   `json:"omit"`     //json 加 omitempty
	Tags     string `json:"tags"`     //自定义Tag
}

type Message struct {
	Name     string    `json:"name"`     //结构体名称
	Fields   []Field   `json:"fields"`   //结构体参数
	Messages []Message `json:"messages"` //请求参数字段结构体
}

type Msg struct {
	Name   string  `json:"name"`   //结构体名称
	Fields []Field `json:"fields"` //结构体参数
}

func (p *Parse) Parse(path, fileType string) (err error) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return
	}

	body, err := io.ReadAll(file)
	if err != nil {
		return
	}

	switch fileType {
	case "toml":
		return p.toml(body)
	case "json":
		err = json.Unmarshal(body, p)
		if err != nil {
			return
		}
	default:
		err = errors.New("not support this file type")
		return
	}

	return
}

func ffn(fields string) (field Field) {
	arrs := strings.Split(fields, ";")

	field = Field{
		Name:     "",
		Class:    "",
		Desc:     "",
		Validate: "",
		From:     "",
		Omit:     false,
		Tags:     "",
	}
	for _, arr := range arrs {
		k, val, _ := strings.Cut(arr, "=")

		switch k {
		case "n":
			field.Name = val
		case "c":
			field.Class = val
		case "d":
			field.Desc = val
		case "v":
			field.Validate = val
		case "f":
			field.From = val
		case "o":
			field.Omit = true
		case "t":
			field.Tags = val
		}
	}

	return
}

func mfn(messages TomlMessage) (message Message) {
	message = Message{
		Name:     messages.Name,
		Fields:   nil,
		Messages: nil,
	}

	for _, fields := range messages.Fields {
		message.Fields = append(message.Fields, ffn(fields))
	}

	for _, msgs := range messages.Messages {
		message.Messages = append(message.Messages, mfn(msgs))
	}
	return
}

func (p *Parse) toml(data []byte) (err error) {
	tp := TomlParse{}

	err = toml.Unmarshal(data, &tp)

	p.Group = tp.Group

	for _, v := range tp.Interfaces {
		one := &Interface{
			Name:           v.Name,
			Desc:           v.Desc,
			Method:         v.Method,
			LastPath:       v.LastPath,
			MidType:        v.MidType,
			ReqContentType: v.ReqContentType,
			ResContentType: v.ResContentType,
			Req: Message{
				Name:     v.Name,
				Fields:   nil,
				Messages: nil,
			},
			Res: Message{
				Name:     v.Name,
				Fields:   nil,
				Messages: nil,
			},
			Msgs: nil,
		}

		for _, fields := range v.Req.Fields {
			one.Req.Fields = append(one.Req.Fields, ffn(fields))
		}

		for _, fields := range v.Res.Fields {
			one.Res.Fields = append(one.Res.Fields, ffn(fields))
		}

		for _, messages := range v.Req.Messages {
			one.Req.Messages = append(one.Req.Messages, mfn(messages))
		}

		for _, messages := range v.Res.Messages {
			one.Res.Messages = append(one.Res.Messages, mfn(messages))
		}

		for _, messages := range v.Msgs {
			msg := Msg{
				Name:   messages.Name,
				Fields: nil,
			}

			for _, fields := range messages.Fields {
				msg.Fields = append(msg.Fields, ffn(fields))
			}

			one.Msgs = append(one.Msgs, msg)
		}

		p.Interfaces = append(p.Interfaces, one)

	}

	return
}

type TomlParse struct {
	Group      Group            `toml:"group"`      //接口组名称
	Interfaces []*TomlInterface `toml:"interfaces"` //接口
}

type TomlInterface struct {
	Name           string      `toml:"name"`             //接口名称
	Desc           string      `toml:"desc"`             //接口备注
	Method         string      `toml:"method"`           //请求方法
	LastPath       string      `toml:"last_path"`        //最后一节路由
	MidType        interface{} `toml:"mid_type"`         //中间件类型
	ReqContentType string      `toml:"req_content_type"` //接参数形式
	ResContentType string      `toml:"res_content_type"` //返回参数形式
	//Header         []Field     `json:"header"`           //请求头参数
	//Uri            []Field     `json:"uri"`              //uri参数
	Req  TomlMessage `toml:"req"`  //请求参数
	Res  TomlMessage `toml:"res"`  //返回数据
	Msgs []TomlMsg   `toml:"msgs"` //通用结构体
}

type TomlMessage struct {
	Name     string        `toml:"name"`     //结构体名称
	Fields   []string      `toml:"fields"`   //结构体参数
	Messages []TomlMessage `toml:"messages"` //请求参数字段结构体
}

type TomlMsg struct {
	Name   string   `toml:"name"`   //结构体名称
	Fields []string `toml:"fields"` //结构体参数
}
