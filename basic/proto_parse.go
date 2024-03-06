package basic

import (
	"encoding/json"
	"errors"
	"github.com/pelletier/go-toml/v2"
	"io"
	"os"
)

type Parse struct {
	Group      Group       `json:"group"`      //接口组名称
	Interfaces []Interface `json:"interfaces"` //接口
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
	Header         []Field     `json:"header"`           //请求头参数
	Uri            []Field     `json:"uri"`              //uri参数
	Req            Message     `json:"req"`              //请求参数
	Res            Message     `json:"res"`              //返回数据
	Msgs           []Msg       `json:"msgs"`             //通用结构体
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
		err = toml.Unmarshal(body, p)
		if err != nil {
			return
		}
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
